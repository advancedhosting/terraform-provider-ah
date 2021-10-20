package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHVolumePlans_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_volume_plans" "test" {
      filter {
		key = "currency"
		values = ["usd"]
	  }
      sort {
        key = "price"
        direction = "asc"
      }
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.slug"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.price"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.currency"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.min_size"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.max_size"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.datacenter_id"),
					resource.TestCheckResourceAttrSet("data.ah_volume_plans.test", "plans.0.datacenter_slug"),
				),
			},
		},
	})
}
