package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAHIPs() *schema.Resource {
	allowedFilterKeys := []string{"id", "reverse_dns"}
	allowedSortingKeys := []string{"id", "created_at", "reverse_dns", "ip_address"}
	return &schema.Resource{
		Read: dataSourceAHIPsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reverse_dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAHIPListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		switch key {
		case "ip_address":
			key = "address"
		}

		sorting := &ah.Sorting{
			Key:   m["key"].(string),
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHIPsListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		filter := &ah.InFilter{
			Keys:   []string{m["key"].(string)},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHIPsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHIPsListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHIPListSorting(v.(*schema.Set))
	}

	ipAddresses, err := client.IPAddresses.List(context.Background(), options)
	if err != nil {
		return err
	}

	if err = dataSourceAHIPsSchema(d, meta, ipAddresses); err != nil {
		return err
	}
	return nil
}

func dataSourceAHIPsSchema(d *schema.ResourceData, meta interface{}, ipAddresses []ah.IPAddress) error {
	ips := make([]map[string]interface{}, len(ipAddresses))
	for i, ipAddress := range ipAddresses {
		ip := map[string]interface{}{
			"id":               ipAddress.ID,
			"ip_address":       ipAddress.Address,
			"datacenter":       ipAddress.DatacenterFullName,
			"type":             ipAddress.Type,
			"reverse_dns":      ipAddress.ReverseDNS,
			"cloud_server_ids": ipAddress.InstanceIDs,
			"created_at":       ipAddress.CreatedAt,
		}
		if primary, err := isPrimaryIP(&ipAddress, meta); err == nil {
			ip["primary"] = primary
		}
		ips[i] = ip
	}
	if err := d.Set("ips", ips); err != nil {
		return fmt.Errorf("unable to set ips attribute: %s", err)
	}
	d.SetId(resource.UniqueId())

	return nil
}

func isPrimaryIP(ipAddress *ah.IPAddress, meta interface{}) (bool, error) {
	client := meta.(*ah.APIClient)
	if ipAddress.Type != "public" {
		return false, fmt.Errorf("IP with type `%s` can not be primary", ipAddress.Type)
	}
	if len(ipAddress.InstanceIDs) != 1 {
		return false, fmt.Errorf("There are no assigned instances with ip %s", ipAddress.Address)
	}

	instanceID := ipAddress.InstanceIDs[0]

	instance, err := client.Instances.Get(context.Background(), instanceID)
	if err != nil {
		return false, err
	}

	instancePrimaryIPAddr, err := instance.PrimaryIPAddr()
	if err != nil {
		return false, err
	}
	return ipAddress.ID == instancePrimaryIPAddr.IPAddressID, nil

}
