package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHVolumeProducts_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_volume_products" "test" {
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.price"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.currency"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.min_size"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.max_size"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.datacenters.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.datacenters.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_volume_products.test", "products.0.datacenters.0.full_name"),
				),
			},
		},
	})
}
