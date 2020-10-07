package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAHDatacenters_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_datacenters" "test" {
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_datacenters.test", "datacenters.#", "2"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.full_name"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.region_id"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.region_name"),
					resource.TestCheckResourceAttrSet("data.ah_datacenters.test", "datacenters.0.region_country_code"),
				),
			},
		},
	})
}
