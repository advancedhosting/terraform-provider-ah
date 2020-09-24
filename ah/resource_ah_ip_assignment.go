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
)

func resourceAHIPAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHIPAssignmentCreate,
		Read:   resourceAHIPAssignmentRead,
		Update: resourceAHIPAssignmentUpdate,
		Delete: resourceAHIPAssignmentDelete,
		Schema: map[string]*schema.Schema{
			"cloud_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAHIPAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	request := &ah.IPAddressAssignmentCreateRequest{
		InstanceID: d.Get("cloud_server_id").(string),
	}

	ipAddressAttr := d.Get("ip_address").(string)

	if _, err := uuid.Parse(ipAddressAttr); err != nil {
		request.IPAddressID = ipAddressAttr
	} else {
		ipAddress, err := ipAddressByIP(ipAddressAttr, meta)
		if err != nil {
			return err
		}
		request.IPAddressID = ipAddress.ID
	}

	ipAssignment, err := client.IPAddressAssignments.Create(context.Background(), request)
	if err != nil {
		return err
	}

	if err := waitForIPAssignmentStatus(d, meta, []string{""}, []string{d.Id()}); err != nil {
		return fmt.Errorf(
			"Error waiting for ip assignment %s: %v", d.Id(), err)
	}

	d.SetId(ipAssignment.ID)

	if d.Get("primary").(bool) {
		if err := setIPAsPrimary(d, meta); err != nil {
			return nil
		}
	}

	return resourceAHIPAssignmentRead(d, meta)

}

func resourceAHIPAssignmentRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ah.APIClient)
	instanceID := d.Get("cloud_server_id").(string)

	instance, err := client.Instances.Get(context.Background(), instanceID)
	if err != nil {
		return err
	}

	d.Set("primary", instance.PrimaryInstanceIPAddressID == d.Id())
	return nil

}

func resourceAHIPAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	if d.HasChange("primary") {
		if d.Get("primary").(bool) {
			if err := setIPAsPrimary(d, meta); err != nil {
				return nil
			}
		}
		d.SetPartial("primary")
	}

	return resourceAHIPAssignmentRead(d, meta)
}

func resourceAHIPAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.IPAddressAssignments.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error deleting ip address assignment (%s): %s", d.Id(), err)
	}

	// if err := waitForIPAssignmentStatus(d, meta, []string{d.Id()}, []string{""}); err != nil {
	// 	return fmt.Errorf(
	// 		"Error waiting for removing ip assignment %s: %v", d.Id(), err)
	// }

	//TODO Be sure assigmnent was removed!!! (See WCS-3540)
	time.Sleep(20 * time.Second)

	return nil
}

func setIPAsPrimary(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ah.APIClient)
	instanceID := d.Get("cloud_server_id").(string)

	action, err := client.Instances.SetPrimaryIP(context.Background(), instanceID, d.Id())

	if err != nil {
		return err
	}

	if err := waitForAction(action.ID, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for setting primary ip %s: %v", d.Id(), err)
	}
	return nil
}

func waitForAction(actionID string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	instanceID := d.Get("cloud_server_id").(string)

	stateRefreshFunc := func() (interface{}, string, error) {
		action, err := client.Instances.ActionInfo(context.Background(), instanceID, actionID)
		if err != nil {
			log.Printf("Error getting action: %v", err)
			return nil, "", err
		}

		return action.ID, action.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    []string{"pending", "running"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"success"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for setting primary ip: %s", err)
	}

	return nil
}

func waitForIPAssignmentStatus(d *schema.ResourceData, meta interface{}, pending, target []string) error {
	client := meta.(*ah.APIClient)
	instanceID := d.Get("cloud_server_id").(string)

	stateRefreshFunc := func() (interface{}, string, error) {
		instance, err := client.Instances.Get(context.Background(), instanceID)
		if err != nil || instance == nil {
			log.Printf("Error getting instance: %v", err)
			return nil, "", err
		}

		for _, ipAddress := range instance.IPAddresses {
			if ipAddress.ID == d.Id() {
				return instance.ID, ipAddress.ID, nil
			}
		}

		return instance.ID, "", nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    pending,
		Refresh:    stateRefreshFunc,
		Target:     target,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance to reach desired assignment status: %s", err)
	}

	return nil
}

func ipAddressByIP(ip string, meta interface{}) (*ah.IPAddress, error) {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{
		Filters: []ah.FilterInterface{
			&ah.ContFilter{
				Keys:  []string{"address"},
				Value: ip,
			},
		},
	}

	ipAddresses, err := client.IPAddresses.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	if len(ipAddresses) != 1 {
		return nil, ah.ErrResourceNotFound
	}
	return &ipAddresses[0], nil
}
