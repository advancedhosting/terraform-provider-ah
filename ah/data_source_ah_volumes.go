package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHVolumes() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "state", "product_id", "size", "file_system", "cloud_server_id"}
	allowedSortingKeys := []string{"id", "name", "state", "product_id", "size", "file_system", "created_at"}
	return &schema.Resource{
		Read: dataSourceAHVolumesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"volumes": {
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
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAHVolumeListSorting(set *schema.Set) []*ah.Sorting {
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

func buildAHVolumesListFilter(set *schema.Set) []ah.FilterInterface {
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
		}

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHVolumesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHVolumesListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHVolumeListSorting(v.(*schema.Set))
	}

	volumes, err := allVolumes(client, options)
	if err != nil {
		return err
	}

	if err = dataSourceAHVolumesSchema(d, meta, volumes); err != nil {
		return err
	}
	return nil
}

func dataSourceAHVolumesSchema(d *schema.ResourceData, meta interface{}, volumes []ah.Volume) error {
	allVolumes := make([]map[string]interface{}, len(volumes))
	var ids string
	for i, volume := range volumes {
		volumeInfo := map[string]interface{}{
			"id":          volume.ID,
			"name":        volume.Name,
			"state":       volume.State,
			"product":     volume.ProductID,
			"size":        volume.Size,
			"file_system": volume.FileSystem,
			"created_at":  volume.CreatedAt,
		}
		if volume.Instance != nil {
			volumeInfo["cloud_server_id"] = volume.Instance.ID
		}
		allVolumes[i] = volumeInfo
		ids += volume.ID
	}
	if err := d.Set("volumes", allVolumes); err != nil {
		return fmt.Errorf("unable to set volumes attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}

func allVolumes(client *ah.APIClient, options *ah.ListOptions) ([]ah.Volume, error) {
	meta := &ah.ListMetaOptions{
		Page: 1,
	}

	options.Meta = meta

	var allVolumes []ah.Volume

	for {
		volumes, meta, err := client.Volumes.List(context.Background(), options)

		if err != nil {
			return nil, fmt.Errorf("Error list volumes: %s", err)
		}

		allVolumes = append(allVolumes, volumes...)
		if meta.IsLastPage() {
			break
		}

		options.Meta.Page++
	}

	return allVolumes, nil
}
