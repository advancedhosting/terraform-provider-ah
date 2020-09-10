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

func resourceAHPrivateNetworkConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHPrivateNetworkConnectionCreate,
		Read:   resourceAHPrivateNetworkConnectionRead,
		Update: resourceAHPrivateNetworkConnectionUpdate,
		Delete: resourceAHPrivateNetworkConnectionDelete,
		Schema: map[string]*schema.Schema{
			"cloud_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAHPrivateNetworkConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	privateNetworkID := d.Get("private_network_id").(string)

	request := &ah.InstancePrivateNetworkCreateRequest{
		InstanceID:       d.Get("cloud_server_id").(string),
		PrivateNetworkID: privateNetworkID,
	}

	if attr, ok := d.GetOk("ip_address"); ok {
		request.IP = attr.(string)
	}

	instancePrivateNetwork, err := client.InstancePrivateNetworks.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating instance private network: %s", err)
	}

	d.SetId(instancePrivateNetwork.ID)
	d.Set("private_network_id", instancePrivateNetwork.PrivateNetwork.ID)

	if err = waitForInstanceConnectionToPrivateNetwork(d, meta); err != nil {
		return err
	}

	return resourceAHPrivateNetworkConnectionRead(d, meta)

}

func resourceAHPrivateNetworkConnectionRead(d *schema.ResourceData, meta interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	instancePrivateNetwork, err := instancePrivateNetworkConnection(privateNetworkID, d, meta)
	if err != nil {
		return err
	}

	d.Set("cloud_server_id", instancePrivateNetwork.Instance.ID)
	d.Set("ip_address", instancePrivateNetwork.IP)

	return nil
}

func resourceAHPrivateNetworkConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	d.Partial(true)

	if d.HasChange("ip_address") {
		updateRequest := &ah.InstancePrivateNetworkUpdateRequest{
			IP: d.Get("ip_address").(string),
		}

		if _, err := client.InstancePrivateNetworks.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf("Error changing ip (%s): %s", d.Id(), err)
		}
		if err := waitForInstanceConnectionToPrivateNetwork(d, meta); err != nil {
			return err
		}
		d.SetPartial("ip_address")
	}

	return resourceAHPrivateNetworkConnectionRead(d, meta)
}

func resourceAHPrivateNetworkConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if _, err := client.InstancePrivateNetworks.Delete(context.Background(), d.Id()); err != nil {
		if err == ah.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting instance private network (%s): %s", d.Id(), err)
	}
	if err := waitForInstancePrivateNetworkDestroy(d, meta); err != nil {
		return err
	}
	return nil
}

func instancePrivateNetworkConnection(privateNetworkID string, d *schema.ResourceData, meta interface{}) (*ah.InstancePrivateNetworkInfo, error) {
	client := meta.(*ah.APIClient)
	privateNetwork, err := client.PrivateNetworks.Get(context.Background(), privateNetworkID)
	if err != nil {
		return nil, err
	}
	for _, instancePrivateNetwork := range privateNetwork.InstancePrivateNetworks {
		if instancePrivateNetwork.ID == d.Id() {
			return &instancePrivateNetwork, nil
		}
	}
	return nil, ah.ErrResourceNotFound
}

func waitForInstanceConnectionToPrivateNetwork(d *schema.ResourceData, meta interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	stateRefreshFunc := func() (interface{}, string, error) {
		instancePrivateNetwork, err := instancePrivateNetworkConnection(privateNetworkID, d, meta)
		if err != nil || instancePrivateNetwork == nil {
			log.Printf("Error on waitForInstanceConnectionToPrivateNetwork: %v", err)
			return nil, "", err
		}
		return instancePrivateNetwork.ID, instancePrivateNetwork.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    []string{"connecting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"connected"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance connection to private network: %s", err)
	}

	return nil
}

func waitForInstancePrivateNetworkDestroy(d *schema.ResourceData, meta interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	stateRefreshFunc := func() (interface{}, string, error) {
		instancePrivateNetwork, err := instancePrivateNetworkConnection(privateNetworkID, d, meta)
		if err == ah.ErrResourceNotFound {
			return d.Id(), "disconnected", nil
		}
		return instancePrivateNetwork.ID, instancePrivateNetwork.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    []string{"disconnecting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"disconnected"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance private network to be destroyed: %s", err)
	}

	return nil
}
