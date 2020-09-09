package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAHPrivateNetworks() *schema.Resource {
	allowedFilterKeys := []string{"id", "ip_range", "name", "cloud_server_id", "state"}
	allowedSortingKeys := []string{"id", "ip_range", "name", "created_at", "state"}
	return &schema.Resource{
		Read: dataSourceAHPrivateNetworksRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"private_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildAHPrivateNetworksListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		sorting := &ah.Sorting{
			Key:   m["key"].(string),
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHPrivateNetworksListFilter(set *schema.Set) []ah.FilterInterface {
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

// TODO Wait for Ransack for private network (WCS-3548)
func dataSourceAHPrivateNetworksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHPrivateNetworksListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHPrivateNetworksListSorting(v.(*schema.Set))
	}

	privateNetworks, err := client.PrivateNetworks.List(context.Background(), options)
	if err != nil {
		return err
	}

	if err = dataSourceAHPrivateNetworksSchema(d, meta, privateNetworks); err != nil {
		return err
	}
	return nil
}

func dataSourceAHPrivateNetworksSchema(d *schema.ResourceData, meta interface{}, privateNetworks []ah.PrivateNetwork) error {
	client := meta.(*ah.APIClient)
	pns := make([]map[string]interface{}, len(privateNetworks))
	for i, privateNetwork := range privateNetworks {
		pn := map[string]interface{}{
			"id":         privateNetwork.ID,
			"ip_range":   privateNetwork.CIDR,
			"name":       privateNetwork.Name,
			"state":      privateNetwork.State,
			"created_at": privateNetwork.CreatedAt,
		}
		privateNetworkInfo, err := client.PrivateNetworks.Get(context.Background(), privateNetwork.ID)
		if err != nil {
			return err
		}
		cloudServers := make([]map[string]interface{}, len(privateNetworkInfo.InstancePrivateNetworks))
		for i, cloudServerInfo := range privateNetworkInfo.InstancePrivateNetworks {
			cloudServer := map[string]interface{}{
				"id": cloudServerInfo.ID,
				"ip": cloudServerInfo.IP,
			}
			cloudServers[i] = cloudServer
		}
		pn["cloud_server_ids"] = cloudServers
		pns[i] = pn
	}

	if err := d.Set("private_networks", pns); err != nil {
		return fmt.Errorf("unable to set private_networks attribute: %s", err)
	}
	d.SetId(resource.UniqueId())

	return nil
}
