package ah

import (
	"context"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataSourceAHVolumePlans() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "slug", "price", "currency", "min_size", "max_size", "datacenter_id", "datacenter_slug"}
	allowedSortingKeys := []string{"id", "name", "slug", "price", "currency", "min_size", "max_size", "datacenter_id", "datacenter_slug"}
	return &schema.Resource{
		Read: dataSourceAHVolumePlansRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAHVolumePlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	VolumePlans, err := allVolumePlans(client)
	if err != nil {
		return err
	}

	if err = dataSourceAHVolumePlansSchema(d, VolumePlans, client); err != nil {
		return err
	}
	return nil
}

func dataSourceAHVolumePlansSchema(d *schema.ResourceData, VolumePlans []ah.VolumePlan, client *ah.APIClient) error {
	var volumePlansData = make([]map[string]interface{}, len(VolumePlans))
	datacenters, err := datacentersInfo(client)
	if err != nil {
		return err
	}

	var ids string
	for i, volumePlan := range VolumePlans {
		datacenterID := volumePlan.CustomAttributes.DatacenterIds[0]
		volumePlanInfo := map[string]interface{}{
			"id":              volumePlan.ID,
			"name":            volumePlan.Name,
			"slug":            volumePlan.CustomAttributes.Slug,
			"currency":        volumePlan.Currency,
			"min_size":        volumePlan.CustomAttributes.MinSize,
			"max_size":        volumePlan.CustomAttributes.MaxSize,
			"datacenter_id":   datacenterID,
			"datacenter_slug": datacenters[datacenterID].Slug,
		}

		for _, price := range volumePlan.Prices {
			if price.Type == "overuse,volume_du" {
				volumePlanInfo["price"] = price.Price
				break
			}
		}
		ids += strconv.Itoa(volumePlan.ID)
		volumePlansData[i] = volumePlanInfo
	}

	volumePlansData = filterPlans(d, volumePlansData)
	sortPlans(d, volumePlansData)

	if err := d.Set("plans", volumePlansData); err != nil {
		return fmt.Errorf("unable to set plans attribute: %s", err)
	}
	d.SetId(generateHash(ids))
	return nil
}

func allVolumePlans(client *ah.APIClient) ([]ah.VolumePlan, error) {

	VolumePlans, err := client.VolumePlans.List(context.Background())

	if err != nil {
		return nil, err
	}

	return VolumePlans, nil
}

func datacentersInfo(client *ah.APIClient) (map[string]ah.Datacenter, error) {
	datacenters, err := client.Datacenters.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	var datacentersMap = make(map[string]ah.Datacenter, len(datacenters))
	for _, dc := range datacenters {
		datacentersMap[dc.ID] = dc
	}
	return datacentersMap, nil
}
