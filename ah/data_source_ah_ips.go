package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHIPs() *schema.Resource {
	allowedFilterKeys := []string{"id", "ip_address", "type", "datacenter", "reverse_dns", "cloud_server_id"}
	allowedSortingKeys := []string{"id", "ip_address", "type", "datacenter", "reverse_dns", "cloud_server_id", "created_at"}
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

func buildAHIPsListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		key := m["key"].(string)

		switch key {
		case "ip_address":
			key = "address"
		case "type":
			key = "address_type"
		case "cloud_server_id":
			key = "instances_id"
		case "datacenter":
			key = "datacenter_id"
		}

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func buildAHIPListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		switch key {
		case "ip_address":
			key = "address"
		case "type":
			key = "ip_network_type"
		case "cloud_server_ids":
			key = "instances_id"
		case "datacenter":
			key = "datacenter_id"
		}

		sorting := &ah.Sorting{
			Key:   key,
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
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
	var ids string
	for i, ipAddress := range ipAddresses {
		ip := map[string]interface{}{
			"id":               ipAddress.ID,
			"ip_address":       ipAddress.Address,
			"datacenter":       ipAddress.DatacenterFullName, // Replace with Datacenter ID after WCS-3498
			"type":             ipAddress.Type,
			"reverse_dns":      ipAddress.ReverseDNS,
			"cloud_server_ids": ipAddress.InstanceIDs,
			"created_at":       ipAddress.CreatedAt,
		}
		if primary, err := isPrimaryIP(&ipAddress, meta); err == nil {
			ip["primary"] = primary
		}
		ips[i] = ip
		ids += ipAddress.ID
	}
	if err := d.Set("ips", ips); err != nil {
		return fmt.Errorf("unable to set ips attribute: %s", err)
	}
	d.SetId(generateHash(ids))

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
