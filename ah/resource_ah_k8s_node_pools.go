package ah

import (
	"context"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var NodePoolSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
	},
	"nodes_count": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"created_at": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"auto_scale": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"min_count": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"max_count": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"public_properties": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
	},
	"private_properties": {
		Type:     schema.TypeMap,
		Optional: true,
	},
	"cluster_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	//"nodes": {
	//	Type:     schema.TypeList,
	//	Optional: true,
	//},
}

func resourceAHK8sNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAHK8sNodePoolCreate,
		ReadContext:   resourceAHK8sNodePoolRead,
		UpdateContext: resourceAHK8sNodePoolUpdate,
		DeleteContext: resourceAHK8sNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: NodePoolSchema,
	}
}

func resourceAHK8sNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	request, err := MakeCreateK8sNodePoolRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if clusterID, ok := d.GetOk("cluster_id"); !ok {
		return diag.Errorf("resource ClusterID is required")
	} else {
		nodePool, err := client.KubernetesClusters.CreateNodePool(ctx, clusterID.(string), request)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(nodePool.ID)
	}

	return resourceAHK8sNodePoolRead(ctx, d, meta)
}

func resourceAHK8sNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	clusterID, nodePoolID := d.Get("cluster_id").(string), d.Id()

	if clusterID == "" || nodePoolID == "" {
		return diag.Errorf("resource ClusterID and ID is required")
	}

	nodePool, err := client.KubernetesClusters.GetNodePool(ctx, clusterID, nodePoolID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", nodePool.Name)
	d.Set("type", nodePool.Type)
	d.Set("nodes_count", nodePool.Count)
	d.Set("created_at", nodePool.CreatedAt)
	d.Set("labels", nodePool.Labels)
	d.Set("auto_scale", nodePool.AutoScale)
	d.Set("min_count", nodePool.MinCount)
	d.Set("max_count", nodePool.MaxCount)

	if nodePool.PublicProperties.PlanID != 0 {
		publicProperties := map[string]int{"plan_id": nodePool.PublicProperties.PlanID}
		d.Set("public_properties", publicProperties)
	} else {
		privateProperties := map[string]interface{}{
			"network_id":      nodePool.PrivateProperties.NetworkID,
			"cluster_id":      nodePool.PrivateProperties.ClusterID,
			"cluster_node_id": nodePool.PrivateProperties.ClusterNodeID,
			"vcpu":            nodePool.PrivateProperties.Vcpu,
			"ram":             nodePool.PrivateProperties.Ram,
			"disk":            nodePool.PrivateProperties.Disk,
		}
		d.Set("private_properties", privateProperties)
	}

	d.Set("cluster_id", clusterID)

	d.SetId(nodePool.ID)

	return nil
}

func resourceAHK8sNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	request := &ah.UpdateKubernetesNodePoolRequest{}
	clusterID, nodePoolID := d.Get("cluster_id").(string), d.Id()

	if clusterID == "" || nodePoolID == "" {
		return diag.Errorf("resource ClusterID and ID is required")
	}

	if !d.HasChanges("nodes_count") && !d.HasChanges("labels") && !d.HasChanges("auto_scale") && !d.HasChanges("max_count") && !d.HasChanges("min_count") {
		return resourceAHK8sClusterRead(ctx, d, meta)
	}

	if c, ok := d.GetOk("nodes_count"); ok {
		request.Count = c.(int)
	}

	if l, ok := d.GetOk("labels"); ok {
		labels := ah.Labels(l.(map[string]string))
		request.Labels = &labels
	}

	if autoScale, ok := d.GetOk("auto_scale"); ok {
		request.AutoScale = autoScale.(bool)
		request.MaxCount = d.Get("max_count").(int)
		request.MinCount = d.Get("min_count").(int)
	}

	err := client.KubernetesClusters.UpdateNodePool(ctx, clusterID, nodePoolID, request)

	if err != nil {
		return diag.Errorf("Error updating K8S Node Pool (%s): %s", d.Id(), err)
	}

	return resourceAHK8sClusterRead(ctx, d, meta)
}

func resourceAHK8sNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	clusterID, nodePoolID := d.Get("cluster_id").(string), d.Id()

	if clusterID == "" || nodePoolID == "" {
		return diag.Errorf("resource ClusterID abd ID is required")
	}

	err := client.KubernetesClusters.DeleteNodePool(ctx, clusterID, nodePoolID, true)

	if err != nil {
		return diag.Errorf("Error delete K8S Node Pool (%s): %s", d.Id(), err)
	}

	return nil
}

func MakeCreateK8sNodePoolRequest(d *schema.ResourceData) (*ah.CreateKubernetesNodePoolRequest, error) {
	nodePoolType := d.Get("type").(string)
	request := &ah.CreateKubernetesNodePoolRequest{Type: nodePoolType}

	if _, ok := d.GetOk("public_properties"); !ok {
		if _, ok := d.GetOk("private_properties"); !ok {
			return nil, fmt.Errorf("must set either public_properties or private_properties")
		}
	} else {
		if _, ok := d.GetOk("private_properties"); ok {
			return nil, fmt.Errorf("cannot set both public_properties and private_properties")
		}
	}

	if l, ok := d.GetOk("labels"); ok {
		labels := ah.Labels(l.(map[string]string))
		request.Labels = &labels
	}

	if autoScale, ok := d.GetOk("auto_scale"); ok {
		request.AutoScale = autoScale.(bool)
		request.MaxCount = d.Get("max_count").(int)
		request.MinCount = d.Get("min_count").(int)
	} else {
		request.Count = d.Get("nodes_count").(int)
	}

	switch nodePoolType {
	case "public":
		request.PublicProperties = &ah.PublicProperties{PlanID: d.Get("public_properties.plan_id").(int)}
	case "private":
		privateProps := d.Get("private_properties").(map[string]interface{})
		request.PrivateProperties = &ah.PrivateProperties{
			NetworkID:     privateProps["network_id"].(string),
			ClusterID:     privateProps["cluster_id"].(string),
			ClusterNodeID: privateProps["cluster_node_id"].(string),
			Vcpu:          privateProps["vcpu"].(int),
			Ram:           privateProps["ram"].(int),
			Disk:          privateProps["disk"].(int),
		}
	}

	return request, nil
}
