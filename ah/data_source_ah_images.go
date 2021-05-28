package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHImages() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "distribution", "version", "architecture", "slug"}
	allowedSortingKeys := []string{"id", "name", "distribution", "version", "architecture", "slug"}
	return &schema.Resource{
		Read: dataSourceAHImagesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"images": {
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
						"distribution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func buildAHImagesListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		sorting := &ah.Sorting{
			Key:   key,
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHImagesListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		key := m["key"].(string)

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHImagesListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHImagesListSorting(v.(*schema.Set))
	}

	images, err := allImages(client, options)
	if err != nil {
		return err
	}

	if err = dataSourceAHImagesSchema(d, meta, images); err != nil {
		return err
	}
	return nil
}

func dataSourceAHImagesSchema(d *schema.ResourceData, meta interface{}, images []ah.Image) error {
	allImages := make([]map[string]interface{}, len(images))
	var ids string
	for i, image := range images {
		imageInfo := map[string]interface{}{
			"id":           image.ID,
			"name":         image.Name,
			"distribution": image.Distribution,
			"version":      image.Version,
			"architecture": image.Architecture,
			"slug":         image.Slug,
		}
		allImages[i] = imageInfo
		ids += image.ID
	}
	if err := d.Set("images", allImages); err != nil {
		return fmt.Errorf("unable to set images attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}

func allImages(client *ah.APIClient, options *ah.ListOptions) ([]ah.Image, error) {
	meta := &ah.ListMetaOptions{
		Page: 1,
	}

	options.Meta = meta

	var images []ah.Image

	for {
		pageImage, meta, err := client.Images.List(context.Background(), options)

		if err != nil {
			return nil, fmt.Errorf("Error list images: %s", err)
		}

		images = append(images, pageImage...)
		if meta.IsLastPage() {
			break
		}

		options.Meta.Page++
	}

	return images, nil
}
