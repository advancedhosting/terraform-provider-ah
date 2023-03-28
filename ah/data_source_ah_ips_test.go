package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHIPs_Basic(t *testing.T) {

	resourcesConfig := fmt.Sprintf(`
	resource "ah_ip" "test" {
	  count = 2
	  type = "public"
	  datacenter = "%s"
	}`, DatacenterID)

	datasourceConfig := `
	data "ah_ips" "test" {
		filter {
			key = "id"
			values = [ah_ip.test.0.id]
		}
		sort {
			key = "created_at"
			direction = "desc"
		}
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_ips.test", "ips.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_ips.test", "ips.0.id", "ah_ip.test.0", "id"),
					resource.TestCheckResourceAttrSet("data.ah_ips.test", "ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.ah_ips.test", "ips.0.datacenter"),
					resource.TestCheckResourceAttrSet("data.ah_ips.test", "ips.0.type"),
					resource.TestCheckResourceAttrSet("data.ah_ips.test", "ips.0.reverse_dns"),
					resource.TestCheckResourceAttrSet("data.ah_ips.test", "ips.0.created_at"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
