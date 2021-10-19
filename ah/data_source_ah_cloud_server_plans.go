package ah

import (
	"context"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataSourceAHCloudServerPlans() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "slug", "price", "currency", "vcpu", "ram", "disk", "available_on_trial"}
	allowedSortingKeys := []string{"id", "name", "slug", "price", "currency", "vcpu", "ram", "disk", "available_on_trial"}
	return &schema.Resource{
		Read: dataSourceAHCloudServerPlansRead,
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
						"vcpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"available_on_trial": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAHCloudServerPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	cloudServerPlans, err := allCloudServerPlans(client)
	if err != nil {
		return err
	}

	if err = dataSourceAHCloudServerPlansSchema(d, cloudServerPlans); err != nil {
		return err
	}
	return nil
}

func dataSourceAHCloudServerPlansSchema(d *schema.ResourceData, cloudServerPlans []ah.InstancePlan) error {
	cloudServerPlansData := make([]map[string]interface{}, len(cloudServerPlans))
	var ids string
	for i, cloudServerPlan := range cloudServerPlans {
		cloudServerPlanInfo := map[string]interface{}{
			"id":                 cloudServerPlan.ID,
			"name":               cloudServerPlan.Name,
			"slug":               cloudServerPlan.CustomAttributes.Slug,
			"currency":           cloudServerPlan.Currency,
			"available_on_trial": cloudServerPlan.CustomAttributes.AvailableOnTrial,
		}

		cloudServerPlanInfo["vcpu"], _ = strconv.Atoi(cloudServerPlan.CustomAttributes.Vcpu)
		cloudServerPlanInfo["ram"], _ = strconv.Atoi(cloudServerPlan.CustomAttributes.RAM)
		cloudServerPlanInfo["disk"], _ = strconv.Atoi(cloudServerPlan.CustomAttributes.Disk)

		for _, price := range cloudServerPlan.Prices {
			if price.Type == "monthly,vps" {
				cloudServerPlanInfo["price"] = price.Price
				break
			}
		}
		ids += strconv.Itoa(cloudServerPlan.ID)
		cloudServerPlansData[i] = cloudServerPlanInfo
	}

	cloudServerPlansData = filterPlans(d, cloudServerPlansData)
	sortPlans(d, cloudServerPlansData)

	if err := d.Set("plans", cloudServerPlansData); err != nil {
		return fmt.Errorf("unable to set plans attribute: %s", err)
	}
	d.SetId(generateHash(ids))
	return nil
}

func allCloudServerPlans(client *ah.APIClient) ([]ah.InstancePlan, error) {

	cloudServerPlans, err := client.InstancePlans.List(context.Background())

	if err != nil {
		return nil, err
	}

	return cloudServerPlans, nil
}
