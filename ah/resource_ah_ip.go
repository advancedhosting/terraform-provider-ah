package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAHIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHIPCreate,
		Read:   resourceAHIPRead,
		Update: resourceAHIPUpdate,
		Delete: resourceAHIPDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "anycast"}, false),
			},
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reverse_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAHIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	addressType := d.Get("type").(string)
	request := &ah.IPAddressCreateRequest{
		Type: addressType,
	}

	if addressType == "public" {
		attr, ok := d.GetOk("datacenter")
		datacenter := attr.(string)
		if !ok || datacenter == "" {
			return fmt.Errorf("Datacenter is required for public ip")
		}
		request.DatacenterID = datacenter
	}

	if attr, ok := d.GetOk("reverse_dns"); ok {
		request.ReverseDNS = attr.(string)
	}

	ipAddress, err := client.IPAddresses.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating ip address: %s", err)
	}

	d.SetId(ipAddress.ID)
	return resourceAHIPRead(d, meta)

}

func resourceAHIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	ipAddress, err := client.IPAddresses.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("reverse_dns", ipAddress.ReverseDNS)
	d.Set("ip_address", ipAddress.Address)
	d.Set("created_at", ipAddress.CreatedAt)

	return nil
}

func resourceAHIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	d.Partial(true)

	if d.HasChange("reverse_dns") {
		reverseDNS := d.Get("reverse_dns").(string)

		updateRequest := &ah.IPAddressUpdateRequest{
			ReverseDNS: reverseDNS,
		}

		if _, err := client.IPAddresses.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"Error changing reverse_dns (%s): %s", d.Id(), err)
		}
		d.SetPartial("reverse_dns")
	}

	return resourceAHIPRead(d, meta)
}

func resourceAHIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.IPAddresses.Delete(context.Background(), d.Id()); err != nil {
		if err == ah.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf(
			"Error deleting ip address (%s): %s", d.Id(), err)
	}
	return nil
}
