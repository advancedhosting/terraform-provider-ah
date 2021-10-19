package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHVolumeProducts() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "slug", "price", "currency", "min_size", "max_size", "datacenter_id", "datacenter_name", "datacenter_slug", "datacenter_full_name"}
	allowedSortingKeys := []string{"id", "name", "slug", "price", "currency", "min_size", "max_size", "datacenter_id", "datacenter_name", "datacenter_slug", "datacenter_full_name"}
	return &schema.Resource{
		Read: dataSourceAHVolumeProductsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"products": {
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
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"currency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildAHVolumeProductsListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		switch key {
		case "cloud_server_id":
			key = "instance_id"
		case "datacenter_id":
			key = "datacenters_id"
		case "datacenter_name":
			key = "datacenters_name"
		case "datacenter_slug":
			key = "datacenters_api_slug"
		case "datacenter_full_name":
			key = "datacenters_full_name"
		}

		sorting := &ah.Sorting{
			Key:   key,
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHVolumeProductsListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		key := m["key"].(string)

		switch key {
		case "cloud_server_id":
			key = "instance_id"
		case "datacenter_id":
			key = "datacenters_id"
		case "datacenter_name":
			key = "datacenters_name"
		case "datacenter_slug":
			key = "datacenters_api_slug"
		case "datacenter_full_name":
			key = "datacenters_full_name"
		}

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHVolumeProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHVolumeProductsListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHVolumeProductsListSorting(v.(*schema.Set))
	}

	volumes, err := allVolumeProducts(client, options)
	if err != nil {
		return err
	}

	if err = dataSourceAHVolumeProductsSchema(d, meta, volumes); err != nil {
		return err
	}
	return nil
}

func dataSourceAHVolumeProductsSchema(d *schema.ResourceData, meta interface{}, volumeProducts []ah.VolumeProduct) error {
	volumeProductsData := make([]map[string]interface{}, len(volumeProducts))
	var ids string
	for i, volumeProduct := range volumeProducts {
		volumeProductInfo := map[string]interface{}{
			"id":       volumeProduct.ID,
			"name":     volumeProduct.Name,
			"slug":     volumeProduct.Slug,
			"price":    volumeProduct.Price,
			"currency": volumeProduct.Currency,
			"min_size": volumeProduct.MinSize,
			"max_size": volumeProduct.MaxSize,
		}

		if len(volumeProduct.DatacenterIDs) > 0 {
			var datacenters []map[string]string
			for _, datacenterID := range volumeProduct.DatacenterIDs {
				datacenter, err := datacenterInfo(datacenterID, meta)
				if err != nil {
					return err
				}
				item := make(map[string]string)
				item["id"] = datacenter.ID
				item["name"] = datacenter.Name
				item["slug"] = datacenter.Slug
				item["full_name"] = datacenter.FullName
				datacenters = append(datacenters, item)
			}
			volumeProductInfo["datacenters"] = datacenters
		}
		volumeProductsData[i] = volumeProductInfo
		ids += volumeProduct.ID
	}
	if err := d.Set("products", volumeProductsData); err != nil {
		return fmt.Errorf("unable to set products attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}

func allVolumeProducts(client *ah.APIClient, options *ah.ListOptions) ([]ah.VolumeProduct, error) {
	meta := &ah.ListMetaOptions{
		Page: 1,
	}

	options.Meta = meta

	var allvolumeProducts []ah.VolumeProduct

	for {
		volumeProducts, meta, err := client.VolumeProducts.List(context.Background(), options)

		if err != nil {
			return nil, fmt.Errorf("Error getting volume product: %s", err)
		}

		allvolumeProducts = append(allvolumeProducts, volumeProducts...)
		if meta.IsLastPage() {
			break
		}

		options.Meta.Page++
	}

	return allvolumeProducts, nil
}

func datacenterInfo(datacenterID string, meta interface{}) (*ah.Datacenter, error) {
	client := meta.(*ah.APIClient)
	datacenter, err := client.Datacenters.Get(context.Background(), datacenterID)
	if err != nil {
		return nil, err
	}
	return datacenter, nil
}
