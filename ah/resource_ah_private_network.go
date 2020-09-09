package ah

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAHPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHPrivateNetworkCreate,
		Read:   resourceAHPrivateNetworkRead,
		Update: resourceAHPrivateNetworkUpdate,
		Delete: resourceAHPrivateNetworkDelete,
		Schema: map[string]*schema.Schema{
			"ip_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAHPrivateNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	request := &ah.PrivateNetworkCreateRequest{
		CIDR: d.Get("ip_range").(string),
	}

	if attr, ok := d.GetOk("name"); ok {
		request.Name = attr.(string)
	}

	privateNetwork, err := client.PrivateNetworks.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating private network: %s", err)
	}

	d.SetId(privateNetwork.ID)

	if err = waitForPrivateNetworkCreate(d, meta); err != nil {
		return err
	}

	return resourceAHPrivateNetworkRead(d, meta)

}

func resourceAHPrivateNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	privateNetwork, err := client.PrivateNetworks.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", privateNetwork.Name)
	d.Set("state", privateNetwork.State)
	d.Set("created_at", privateNetwork.CreatedAt)

	return nil
}

func resourceAHPrivateNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	d.Partial(true)

	if d.HasChange("ip_range") {
		updateRequest := &ah.PrivateNetworkUpdateRequest{
			CIDR: d.Get("ip_range").(string),
		}

		if _, err := client.PrivateNetworks.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"Error changing ip range (%s): %s", d.Id(), err)
		}
		d.SetPartial("ip_range")
	}

	if d.HasChange("name") {
		updateRequest := &ah.PrivateNetworkUpdateRequest{
			Name: d.Get("name").(string),
		}

		if _, err := client.PrivateNetworks.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"Error changing private network's name (%s): %s", d.Id(), err)
		}
		d.SetPartial("name")
	}

	return resourceAHPrivateNetworkRead(d, meta)
}

func resourceAHPrivateNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.PrivateNetworks.Delete(context.Background(), d.Id()); err != nil {
		if err == ah.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting private network (%s): %s", d.Id(), err)
	}
	if err := waitForPrivateNetworkDestroy(d, meta); err != nil {
		return err
	}
	return nil
}

func waitForPrivateNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		privateNetwork, err := client.PrivateNetworks.Get(context.Background(), d.Id())
		if err != nil || privateNetwork == nil {
			log.Printf("Error on waitForPrivateNetworkCreate: %v", err)
			return nil, "", err
		}
		return privateNetwork.ID, privateNetwork.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    []string{"updating"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"active"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for private network to be active: %s", err)
	}

	return nil
}

func waitForPrivateNetworkDestroy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		privateNetwork, err := client.PrivateNetworks.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		return privateNetwork.ID, privateNetwork.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    []string{"deleting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for private network to be destroyed: %s", err)
	}

	return nil
}
