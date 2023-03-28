package ah

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"private_cloud": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vcpu": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAHK8sClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	request := &ah.ClusterCreateRequest{
		Name:  d.Get("name").(string),
		Count: d.Get("nodes_count").(int),
	}

	datacenterAttr := d.Get("datacenter").(string)
	if _, err := uuid.Parse(datacenterAttr); err != nil {
		datacenterID, err := datacenterIDBySlug(ctx, client, d.Get("datacenter").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.DatacenterID = datacenterID
	} else {
		request.DatacenterID = datacenterAttr
	}

	request.PrivateCloud = d.Get("private_cloud").(bool)
	if request.PrivateCloud {
		request.Vcpu = d.Get("vcpu").(int)
		request.Ram = d.Get("ram").(int)
		request.Disk = d.Get("disk").(int)
	} else {
		planAttr := d.Get("plan").(string)
		if planID, err := strconv.Atoi(planAttr); err == nil {
			request.PlanId = planID
		}
	}

	cluster, err := client.Clusters.Create(ctx, request)

	if err != nil {
		return diag.Errorf("Error creating k8s cluster: %s", err)
	}
	d.SetId(cluster.ID)

	if err := waitForK8sClusterStatus(ctx, []string{"creating"}, []string{"active"}, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for k8s cluster (%s) to become ready: %s", d.Id(), err)
	}

	return resourceAHK8sClusterRead(ctx, d, meta)

}

func resourceAHK8sClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	cluster, err := client.Clusters.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", cluster.Name)
	d.Set("state", cluster.State)
	d.Set("created_at", cluster.CreatedAt)
	d.Set("number", cluster.Number)
	d.Set("nodes_count", cluster.Count)

	return nil

}

func resourceAHK8sClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	if d.HasChange("name") {
		request := &ah.ClusterUpdateRequest{
			Name: d.Get("name").(string),
		}

		err := client.Clusters.Update(ctx, d.Id(), request)

		if err != nil {
			return diag.Errorf(
				"Error renaming k8s cluster (%s): %s", d.Id(), err)
		}
	}

	return resourceAHK8sClusterRead(ctx, d, meta)
}

func resourceAHK8sClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	if err := client.Clusters.Delete(ctx, d.Id()); err != nil {
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
		cluster, err := client.Clusters.Get(context.Background(), d.Id())
		if err != nil {
			log.Printf("Error on waitForK8sClusterStatus: %v", err)
			return nil, "", err
		}
		return cluster.ID, cluster.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      20 * time.Second,
		Pending:    pendingStatuses,
		Refresh:    stateRefreshFunc,
		Target:     targetStatuses,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf(
			"error waiting for k8s cluster to reach desired status %s: %s", targetStatuses, err)
	}

	return nil
}

func waitForK8sClusterDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		cluster, err := client.Clusters.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil {
			log.Printf("Error on waitForK8sClusterDestroy: %v", err)
			return nil, "", err
		}

		return cluster.ID, cluster.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
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
