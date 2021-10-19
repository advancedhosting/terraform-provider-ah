package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHCloudServerPlans_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_cloud_server_plans" "test" {
      filter {
		key = "vcpu"
		values = [1]
	  }
      sort {
        key = "ram"
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
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.slug"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.price"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.currency"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.vcpu"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.ram"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.disk"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_plans.test", "plans.0.available_on_trial"),
				),
			},
		},
	})
}
