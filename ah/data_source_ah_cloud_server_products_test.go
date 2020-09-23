package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAHCloudServerProducts_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_cloud_server_products" "test" {
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.name"),
					//resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.slug"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.price"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.currency"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.vcpu"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.ram"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.disk"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_server_products.test", "products.0.available_on_trial"),
				),
			},
		},
	})
}
