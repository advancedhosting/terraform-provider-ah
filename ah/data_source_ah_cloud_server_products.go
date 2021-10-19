package ah

import (
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//
//import (
//	"context"
//	"fmt"
//
//	"github.com/advancedhosting/advancedhosting-api-go/ah"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//func dataSourceAHCloudServerProducts() *schema.Resource {
//	allowedFilterKeys := []string{"id", "name", "slug", "price", "currency", "vcpu", "ram", "disk", "available_on_trial"}
//	allowedSortingKeys := []string{"id", "name", "slug", "price", "currency", "vcpu", "ram", "disk", "available_on_trial"}
//	return &schema.Resource{
//		Read: dataSourceAHCloudServerProductsRead,
//		Schema: map[string]*schema.Schema{
//			"filter": dataSourceFilterSchema(allowedFilterKeys),
//			"sort":   dataSourceSortingSchema(allowedSortingKeys),
//			"products": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"id": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"name": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"slug": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"price": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"currency": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"vcpu": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"ram": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"disk": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"available_on_trial": {
//							Type:     schema.TypeBool,
//							Computed: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func buildAHCloudServerProductsListSorting(set *schema.Set) []*ah.Sorting {
//	var sortings []*ah.Sorting
//	for _, v := range set.List() {
//		m := v.(map[string]interface{})
//
//		key := m["key"].(string)
//
//		sorting := &ah.Sorting{
//			Key:   key,
//			Order: m["direction"].(string),
//		}
//
//		sortings = append(sortings, sorting)
//	}
//	return sortings
//}
//
func buildAHCloudServerProductsListFilter(set *schema.Set) []ah.FilterInterface {
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

//
//func dataSourceAHCloudServerProductsRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*ah.APIClient)
//	options := &ah.ListOptions{}
//
//	if v, ok := d.GetOk("filter"); ok {
//		options.Filters = buildAHCloudServerProductsListFilter(v.(*schema.Set))
//	}
//
//	if v, ok := d.GetOk("sort"); ok {
//		options.Sortings = buildAHCloudServerProductsListSorting(v.(*schema.Set))
//	}
//
//	cloudServerProducts, err := allCloudServerProducts(client, options)
//	if err != nil {
//		return err
//	}
//
//	if err = dataSourceAHCloudServerProductsSchema(d, meta, cloudServerProducts); err != nil {
//		return err
//	}
//	return nil
//}
//
//func dataSourceAHCloudServerProductsSchema(d *schema.ResourceData, meta interface{}, cloudServerProducts []ah.InstanceProduct) error {
//	cloudServerProductsData := make([]map[string]interface{}, len(cloudServerProducts))
//	var ids string
//	for i, cloudServerProduct := range cloudServerProducts {
//		cloudServerProductInfo := map[string]interface{}{
//			"id":                 cloudServerProduct.ID,
//			"name":               cloudServerProduct.Name,
//			"slug":               cloudServerProduct.Slug,
//			"price":              cloudServerProduct.Price,
//			"currency":           cloudServerProduct.Currency,
//			"vcpu":               cloudServerProduct.Vcpu,
//			"ram":                cloudServerProduct.RAM,
//			"disk":               cloudServerProduct.Disk,
//			"available_on_trial": cloudServerProduct.AvailableOnTrial,
//		}
//		ids += cloudServerProduct.ID
//
//		cloudServerProductsData[i] = cloudServerProductInfo
//	}
//	if err := d.Set("products", cloudServerProductsData); err != nil {
//		return fmt.Errorf("unable to set products attribute: %s", err)
//	}
//	d.SetId(generateHash(ids))
//	return nil
//}
//
//func allCloudServerProducts(client *ah.APIClient, options *ah.ListOptions) ([]ah.InstanceProduct, error) {
//	meta := &ah.ListMetaOptions{
//		Page: 1,
//	}
//
//	options.Meta = meta
//
//	var cloudServerProducts []ah.InstanceProduct
//
//	for {
//		cloudServerProductsPage, meta, err := client.InstanceProducts.List(context.Background(), options)
//
//		if err != nil {
//			return nil, fmt.Errorf("Error list cloud server products: %s", err)
//		}
//
//		cloudServerProducts = append(cloudServerProducts, cloudServerProductsPage...)
//		if meta.IsLastPage() {
//			break
//		}
//
//		options.Meta.Page++
//	}
//
//	return cloudServerProducts, nil
//}
