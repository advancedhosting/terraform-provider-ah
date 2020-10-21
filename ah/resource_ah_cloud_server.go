package ah

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAHCloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHCloudServerCreate,
		Read:   resourceAHCloudServerRead,
		Update: resourceAHCloudServerUpdate,
		Delete: resourceAHCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAHCloudServerImport,
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAHCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
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

	productAttr := d.Get("product").(string)
	if _, err := uuid.Parse(productAttr); err != nil {
		request.ProductSlug = productAttr
	} else {
		request.ProductID = productAttr
	}

	if attr, ok := d.GetOk("ssh_keys"); ok {
		var sshKeyIDs []string
		for _, v := range attr.([]interface{}) {
			if IsUUID(v.(string)) {
				sshKeyIDs = append(sshKeyIDs, v.(string))
			} else {
				sshKey, err := sshKeyByFingerprint(v.(string), meta)
				if err != nil {
					return fmt.Errorf("Error searching ssh key by fingerprint %s: %v", v.(string), err)
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

	instance, err := client.Instances.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	d.SetId(instance.ID)
	if err = waitForStatus([]string{"creating", "stopped"}, []string{"running"}, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for cloud server (%s) to become ready: %s", d.Id(), err)
	}

	//Wait for instance to completely load
	time.Sleep(20 * time.Second)

	return resourceAHCloudServerRead(d, meta)

}

func resourceAHCloudServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ah.APIClient)
	_, err := client.Instances.Get(context.Background(), d.Id())
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceAHCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	instance, err := client.Instances.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", instance.Name)
	d.Set("created_at", instance.CreatedAt)
	d.Set("state", instance.State)
	d.Set("vcpu", instance.Vcpu)
	d.Set("ram", instance.RAM)
	d.Set("disk", instance.RAM)

	var ips []map[string]interface{}
	for _, instanceIPAddress := range instance.IPAddresses {
		item := make(map[string]interface{})
		item["assignment_id"] = instanceIPAddress.ID
		item["ip_address"] = instanceIPAddress.Address
		item["primary"] = instance.PrimaryInstanceIPAddressID == instanceIPAddress.ID
		ipAddress, err := client.IPAddresses.Get(context.Background(), instanceIPAddress.IPAddressID)
		if err != nil {
			return err
		}
		item["type"] = ipAddress.Type
		item["reverse_dns"] = ipAddress.ReverseDNS
		ips = append(ips, item)

	}
	d.Set("ips", ips)

	return nil

}

func resourceAHCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		client := meta.(*ah.APIClient)

		_, err := client.Instances.Rename(context.Background(), d.Id(), newName)

		if err != nil {
			return fmt.Errorf(
				"Error renaming cloud server (%s): %s", d.Id(), err)
		}
		d.SetPartial("name")
	}

	if d.HasChange("product") {
		client := meta.(*ah.APIClient)

		request := &ah.InstanceUpgradeRequest{}

		newProductAttr := d.Get("product").(string)
		if _, err := uuid.Parse(newProductAttr); err != nil {
			request.ProductSlug = newProductAttr
		} else {
			request.ProductID = newProductAttr
		}

		if err := client.Instances.Upgrade(context.Background(), d.Id(), request); err != nil {
			return fmt.Errorf(
				"Error upgrade instance (%s): %s", d.Id(), err)
		}
		if err := waitForStatus([]string{"updating"}, []string{"running"}, d, meta); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become stopped: %s", d.Id(), err)
		}
		d.SetPartial("product")
	}

	return resourceAHCloudServerRead(d, meta)
}

func resourceAHCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.Instances.PowerOff(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error power_off instance (%s): %s", d.Id(), err)
	}

	if err := waitForStatus([]string{"running"}, []string{"stopped"}, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become stopped: %s", d.Id(), err)
	}

	if err := client.Instances.Destroy(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error destroy instance (%s): %s", d.Id(), err)
	}

	if err := waitForDestroy(d, meta); err != nil {
		return fmt.Errorf(
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
	_, err := stateChangeConf.WaitForState()

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
