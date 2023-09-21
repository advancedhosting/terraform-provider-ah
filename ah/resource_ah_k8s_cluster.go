package ah

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	state "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAHK8sCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAHK8sClusterCreate,
		ReadContext:   resourceAHK8sClusterRead,
		UpdateContext: resourceAHK8sClusterUpdate,
		DeleteContext: resourceAHK8sClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"datacenter": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"private_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"k8s_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"node_pools": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: NodePoolSchema,
				},
			},
		},
	}
}

func resourceAHK8sClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	var datacenterID string
	var nodePools []ah.CreateKubernetesNodePoolRequest

	datacenterAttr := d.Get("datacenter").(string)

	if _, err := uuid.Parse(datacenterAttr); err == nil {
		datacenterID = datacenterAttr
	} else {
		id, err := datacenterIDBySlug(ctx, client, datacenterAttr)
		if err != nil {
			return diag.FromErr(err)
		}
		datacenterID = id
	}

	k8sVersion, err := kubernetesVersion(ctx, client, d.Get("k8s_version").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, nodePool := range d.Get("node_pools").([]interface{}) {
		nodePoolRequest, err := expandCreateNodePoolRequest(nodePool)
		if err != nil {
			return diag.FromErr(err)
		}
		nodePools = append(nodePools, *nodePoolRequest)
	}

	request := &ah.KubernetesClusterCreateRequest{
		Name:         d.Get("name").(string),
		DatacenterID: datacenterID,
		K8sVersion:   k8sVersion,
		NodePools:    nodePools,
	}

	cluster, err := client.KubernetesClusters.Create(ctx, request)

	if err != nil {
		return diag.Errorf("Error creating k8s cluster: %s", err)
	}

	d.SetId(cluster.ID)

	if err := waitForK8sClusterStatus(ctx, []string{"creating", "creation_failed"}, []string{"active"}, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for k8s cluster (%s) to become ready: %s", d.Id(), err)
	}

	return resourceAHK8sClusterRead(ctx, d, meta)

}

func resourceAHK8sClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	clusterID := d.Id()
	if clusterID == "" {
		return diag.Errorf("resource ID is required")
	}

	cluster, err := client.KubernetesClusters.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", cluster.Name)
	d.Set("datacenter", cluster.DatacenterSlug)
	d.Set("private_network", cluster.PrivateNetworkName)
	d.Set("state", cluster.State)
	d.Set("created_at", cluster.CreatedAt)
	d.Set("number", cluster.Number)
	d.Set("account_id", cluster.AccountID)
	d.Set("k8s_version", cluster.K8sVersion)

	if err = dataSourceAHNodePoolSchema(d, cluster.NodePools); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.ID)

	return nil
}

func resourceAHK8sClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	if !d.HasChanges("name") {
		return resourceAHK8sClusterRead(ctx, d, meta)
	}

	request := &ah.KubernetesClusterUpdateRequest{
		Name: d.Get("name").(string),
	}

	err := client.KubernetesClusters.Update(ctx, d.Id(), request)

	if err != nil {
		return diag.Errorf("Error renaming k8s cluster (%s): %s", d.Id(), err)
	}

	return resourceAHK8sClusterRead(ctx, d, meta)
}

func resourceAHK8sClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	if err := client.KubernetesClusters.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf(
			"Error deleting k8s cluster (%s): %s", d.Id(), err)
	}

	if err := waitForK8sClusterDestroy(ctx, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for k8s cluster (%s) to become deleted: %s", d.Id(), err)
	}

	return nil
}

func waitForK8sClusterStatus(ctx context.Context, pendingStatuses, targetStatuses []string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		cluster, err := client.KubernetesClusters.Get(context.Background(), d.Id())
		if err != nil {
			log.Printf("Error on waitForK8sClusterStatus: %v", err)
			return nil, "", err
		}
		return cluster.ID, cluster.State, nil
	}

	stateChangeConf := state.StateChangeConf{
		Delay:      20 * time.Second,
		Pending:    pendingStatuses,
		Refresh:    stateRefreshFunc,
		Target:     targetStatuses,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf("error waiting for k8s cluster to reach desired status %s: %s", targetStatuses, err)
	}

	return nil
}

func waitForK8sClusterDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		cluster, err := client.KubernetesClusters.Get(context.Background(), d.Id())
		if errors.Is(err, ah.ErrResourceNotFound) {
			return d.Id(), "deleted", nil
		}
		if err != nil {
			log.Printf("Error on waitForK8sClusterDestroy: %v", err)
			return nil, "", err
		}

		return cluster.ID, cluster.State, nil
	}

	stateChangeConf := state.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"active", "deleting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf(
			"error waiting for k8s cluster to reach desired status deleted: %s", err)
	}

	return nil
}

func expandCreateNodePoolRequest(np interface{}) (*ah.CreateKubernetesNodePoolRequest, error) {
	labels := ah.Labels{}
	nodePool := np.(map[string]interface{})

	nodePoolRequest := &ah.CreateKubernetesNodePoolRequest{Type: nodePool["type"].(string)}

	if l, ok := nodePool["labels"].(map[string]interface{}); ok {
		for k, v := range l {
			labels[k] = fmt.Sprintf("%v", v)
		}
		nodePoolRequest.Labels = &labels
	}

	if autoScale, ok := nodePool["auto_scale"]; !ok {
		nodePoolRequest.AutoScale = autoScale.(bool)
		nodePoolRequest.MaxCount = nodePool["max_count"].(int)
		nodePoolRequest.MinCount = nodePool["min_count"].(int)
	} else {
		nodePoolRequest.Count = nodePool["nodes_count"].(int)
	}

	if _, ok := nodePool["public_properties"]; !ok {
		if _, ok := nodePool["private_properties"]; !ok {
			return nil, fmt.Errorf("must set either public_properties or private_properties")
		}
		privateProperties := nodePool["private_properties"].(map[string]interface{})
		nodePoolRequest.PrivateProperties = &ah.PrivateProperties{
			NetworkID:     privateProperties["network_id"].(string),
			ClusterID:     privateProperties["cluster_id"].(string),
			ClusterNodeID: privateProperties["cluster_node_id"].(string),
			Vcpu:          privateProperties["vcpu"].(int),
			Ram:           privateProperties["ram"].(int),
			Disk:          privateProperties["disk"].(int),
		}
	} else {
		planID := nodePool["public_properties"].(map[string]interface{})["plan_id"].(int)
		nodePoolRequest.PublicProperties = &ah.PublicProperties{PlanID: planID}
	}

	return nodePoolRequest, nil
}

func dataSourceAHNodePoolSchema(d *schema.ResourceData, nodePools []ah.KubernetesNodePool) error {
	allNodePools := make([]map[string]interface{}, len(nodePools))
	var ids string
	for i, nodePool := range nodePools {
		nodePoolInfo := map[string]interface{}{
			"id":          nodePool.ID,
			"name":        nodePool.Name,
			"type":        nodePool.Type,
			"nodes_count": nodePool.Count,
			"auto_scale":  nodePool.AutoScale,
			"min_count":   nodePool.MinCount,
			"max_count":   nodePool.MaxCount,
			"labels":      nodePool.Labels,
		}

		if nodePool.PublicProperties.PlanID != 0 {
			nodePoolInfo["public_properties"] = map[string]int{
				"plan_id": nodePool.PublicProperties.PlanID,
			}
		} else {
			nodePoolInfo["private_properties"] = map[string]interface{}{
				"network_id":      nodePool.PrivateProperties.NetworkID,
				"cluster_id":      nodePool.PrivateProperties.ClusterID,
				"cluster_node_id": nodePool.PrivateProperties.ClusterNodeID,
				"vcpu":            nodePool.PrivateProperties.Vcpu,
				"ram":             nodePool.PrivateProperties.Ram,
				"disk":            nodePool.PrivateProperties.Disk,
			}
		}

		allNodePools[i] = nodePoolInfo
		ids += nodePool.ID
	}
	if err := d.Set("node_pools", allNodePools); err != nil {
		return fmt.Errorf("unable to set Node Pools attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}
