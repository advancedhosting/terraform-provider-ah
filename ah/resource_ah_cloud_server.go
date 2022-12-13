package ah

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAHCloudServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAHCloudServerCreate,
		ReadContext:   resourceAHCloudServerRead,
		UpdateContext: resourceAHCloudServerUpdate,
		DeleteContext: resourceAHCloudServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAHCloudServerImport,
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
			"image": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"product": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use plan instead",
			},
			"plan": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"product"},
			},
			"use_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"reverse_dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"private_cloud": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAHCloudServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	request := &ah.InstanceCreateRequest{
		Name:                  d.Get("name").(string),
		CreatePublicIPAddress: d.Get("create_public_ip_address").(bool),
		UseSSHPassword:        d.Get("use_password").(bool),
	}

	datacenterAttr := d.Get("datacenter").(string)
	if _, err := uuid.Parse(datacenterAttr); err != nil {
		request.DatacenterSlug = datacenterAttr
	} else {
		request.DatacenterID = datacenterAttr
	}

	imageAttr := d.Get("image").(string)
	if _, err := uuid.Parse(imageAttr); err != nil {
		request.ImageSlug = imageAttr
	} else {
		request.ImageID = imageAttr
	}

	privateCloud, privateCloudOk := d.GetOk("private_cloud")
	if privateCloudOk {
		request.PrivateCloud = privateCloud.(bool)
		request.ClusterID = d.Get("cluster_id").(string)
		request.NodeID = d.Get("node_id").(string)
		networkID, networkIDOk := d.GetOk("network_id")
		if networkIDOk {
			request.IPNetworkID = networkID.(string)
		}
		request.Vcpu = d.Get("vcpu").(int)
		request.Ram = d.Get("ram").(int)
		request.Disk = d.Get("disk").(int)
	} else {
		var planAttr string
		plan, planOk := d.GetOk("plan")
		product, productOk := d.GetOk("product")
		if !planOk && !productOk {
			return diag.Errorf("one of plan or product must be configured")
		}

		if planOk {
			planAttr = plan.(string)
		} else {
			planAttr = product.(string)
		}

		if planID, err := strconv.Atoi(planAttr); err != nil {
			request.PlanSlug = planAttr
		} else {
			request.PlanID = planID
		}
	}

	if attr, ok := d.GetOk("ssh_keys"); ok {
		var sshKeyIDs []string
		for _, v := range attr.([]interface{}) {
			if IsUUID(v.(string)) {
				sshKeyIDs = append(sshKeyIDs, v.(string))
			} else {
				sshKey, err := sshKeyByFingerprint(v.(string), meta)
				if err != nil {
					return diag.Errorf("Error searching ssh key by fingerprint %s: %v", v.(string), err)
				}
				sshKeyIDs = append(sshKeyIDs, sshKey.ID)

			}

		}
		request.SSHKeyIDs = sshKeyIDs
	}

	if attr, ok := d.GetOk("backups"); ok {
		if attr.(bool) {
			request.SnapshotBySchedule = true
			request.SnapshotPeriod = "weekly"
		}
	}

	instance, err := client.Instances.Create(ctx, request)

	if err != nil {
		return diag.Errorf("Error creating instance: %s", err)
	}

	d.SetId(instance.ID)
	if err = waitForStatus([]string{"creating", "stopped"}, []string{"running"}, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for cloud server (%s) to become ready: %s", d.Id(), err)
	}

	//Wait for instance to completely load
	time.Sleep(20 * time.Second)

	return resourceAHCloudServerRead(ctx, d, meta)

}

func resourceAHCloudServerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ah.APIClient)
	instance, err := client.Instances.Get(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	if instance.Image.Slug != "" {
		d.Set("image", instance.Image.Slug)
	} else {
		d.Set("image", instance.Image.ID)
	}

	if instance.Datacenter.Slug != "" {
		d.Set("datacenter", instance.Datacenter.Slug)
	} else {
		d.Set("datacenter", instance.Datacenter.ID)
	}

	d.Set("product", instance.ProductID)

	sshKeys := make([]string, len(instance.SSHKeys))
	for i, sshKey := range instance.SSHKeys {
		sshKeys[i] = sshKey.ID
	}

	d.Set("ssh_keys", sshKeys)

	return []*schema.ResourceData{d}, nil
}

func resourceAHCloudServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	instance, err := client.Instances.Get(context.Background(), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", instance.Name)
	d.Set("created_at", instance.CreatedAt)
	d.Set("state", instance.State)
	d.Set("vcpu", instance.Vcpu)
	d.Set("ram", instance.RAM)
	d.Set("disk", instance.Disk)
	d.Set("backups", instance.SnapshotBySchedule)
	d.Set("use_password", instance.UseSSHPassword)

	var ips []map[string]interface{}
	for _, instanceIPAddress := range instance.IPAddresses {
		item := make(map[string]interface{})
		item["assignment_id"] = instanceIPAddress.ID
		item["ip_address"] = instanceIPAddress.Address
		item["primary"] = instance.PrimaryInstanceIPAddressID == instanceIPAddress.ID
		ipAddress, err := client.IPAddresses.Get(context.Background(), instanceIPAddress.IPAddressID)
		if err != nil {
			return diag.FromErr(err)
		}
		item["type"] = ipAddress.Type
		item["reverse_dns"] = ipAddress.ReverseDNS
		ips = append(ips, item)

	}
	d.Set("ips", ips)

	return nil

}

func resourceAHCloudServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		client := meta.(*ah.APIClient)

		_, err := client.Instances.Rename(ctx, d.Id(), newName)

		if err != nil {
			return diag.Errorf(
				"Error renaming cloud server (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("plan") {
		client := meta.(*ah.APIClient)

		request := &ah.InstanceUpgradeRequest{}

		planAttr := d.Get("plan").(string)
		if planID, err := strconv.Atoi(planAttr); err != nil {
			request.PlanSlug = planAttr
		} else {
			request.PlanID = planID
		}

		if err := client.Instances.Upgrade(ctx, d.Id(), request); err != nil {
			return diag.Errorf(
				"Error upgrade instance (%s): %s", d.Id(), err)
		}
		if err := waitForStatus([]string{"updating"}, []string{"running"}, d, meta); err != nil {
			return diag.Errorf(
				"Error waiting for instance (%s) to become stopped: %s", d.Id(), err)
		}
	}

	return resourceAHCloudServerRead(ctx, d, meta)
}

func resourceAHCloudServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	if err := client.Instances.PowerOff(ctx, d.Id()); err != nil {
		return diag.Errorf(
			"Error power_off instance (%s): %s", d.Id(), err)
	}

	if err := waitForStatus([]string{"running"}, []string{"stopped"}, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become stopped: %s", d.Id(), err)
	}

	if err := client.Instances.Destroy(ctx, d.Id()); err != nil {
		return diag.Errorf(
			"Error destroy instance (%s): %s", d.Id(), err)
	}

	if err := waitForDestroy(d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become deleted: %s", d.Id(), err)
	}

	return nil
}

func waitForStatus(pendingStatuses, targetStatuses []string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		instance, err := client.Instances.Get(context.Background(), d.Id())
		if err != nil || instance == nil {
			log.Printf("Error on InstanceStateRefresh: %v", err)
			return nil, "", err
		}
		return instance.ID, instance.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      20 * time.Second,
		Pending:    pendingStatuses,
		Refresh:    stateRefreshFunc,
		Target:     targetStatuses,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(context.Background())

	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance to reach desired status %s: %s", targetStatuses, err)
	}

	return nil
}

func waitForDestroy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		instance, err := client.Instances.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil || instance == nil {
			log.Printf("Error on waitOnDestory: %v", err)
			return nil, "", err
		}

		return instance.ID, instance.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"stopped", "destroying"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance to reach desired status deleted: %s", err)
	}

	return nil
}

func sshKeyByFingerprint(fingerprint string, meta interface{}) (*ah.SSHKey, error) {
	client := meta.(*ah.APIClient)
	sshKeys, err := allSSHKeysInfo(client, &ah.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, sshKey := range sshKeys {
		if sshKey.Fingerprint == fingerprint {
			return &sshKey, nil
		}
	}
	return nil, ah.ErrResourceNotFound
}
