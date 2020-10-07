package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAHDatacenters() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "slug", "full_name", "region_id", "region_name", "region_country_code"}
	allowedSortingKeys := []string{"id", "name", "slug", "full_name", "region_id", "region_name", "region_country_code"}
	return &schema.Resource{
		Read: dataSourceAHDatacentersRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"datacenters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAHDatacentersListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		switch key {
		case "slug":
			key = "datacenter_slug"
		}

		sorting := &ah.Sorting{
			Key:   key,
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHDatacentersListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		key := m["key"].(string)

		switch key {
		case "slug":
			key = "datacenter_slug"
		}

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHDatacentersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHDatacentersListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHDatacentersListSorting(v.(*schema.Set))
	}

	datacenters, err := client.Datacenters.List(context.Background(), options)
	if err != nil {
		return err
	}

	if err = dataSourceAHDatacentersSchema(d, meta, datacenters); err != nil {
		return err
	}
	return nil
}

func dataSourceAHDatacentersSchema(d *schema.ResourceData, meta interface{}, datacenters []ah.Datacenter) error {
	allDatacenters := make([]map[string]interface{}, len(datacenters))
	for i, datacenter := range datacenters {
		datacenterInfo := map[string]interface{}{
			"id":                  datacenter.ID,
			"name":                datacenter.Name,
			"slug":                datacenter.Slug,
			"full_name":           datacenter.FullName,
			"region_id":           datacenter.Region.ID,
			"region_name":         datacenter.Region.Name,
			"region_country_code": datacenter.Region.CountryCode,
		}

		allDatacenters[i] = datacenterInfo
	}
	if err := d.Set("datacenters", allDatacenters); err != nil {
		return fmt.Errorf("unable to set datacenters attribute: %s", err)
	}
	d.SetId(resource.UniqueId())

	return nil
}
