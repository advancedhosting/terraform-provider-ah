package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHPrivateNetworks_Basic(t *testing.T) {

	resourcesConfig := fmt.Sprintf(`
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}
	
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "ams1"
	  image = "%s"
	  product = "%s"
	}
	
	resource "ah_private_network_connection" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  private_network_id = ah_private_network.test.id
	}`, ImageName, VpsPlanName)

	datasourceConfig := `
	data "ah_private_networks" "test" {
		filter {
			key = "id"
			values = [ah_private_network.test.id]
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
					resource.TestCheckResourceAttr("data.ah_private_networks.test", "private_networks.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_private_networks.test", "private_networks.0.id", "ah_private_network.test", "id"),
					resource.TestCheckResourceAttrSet("data.ah_private_networks.test", "private_networks.0.ip_range"),
					resource.TestCheckResourceAttrSet("data.ah_private_networks.test", "private_networks.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_private_networks.test", "private_networks.0.state"),
					resource.TestCheckResourceAttrSet("data.ah_private_networks.test", "private_networks.0.created_at"),
					resource.TestCheckResourceAttr("data.ah_private_networks.test", "private_networks.0.cloud_servers.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_private_networks.test", "private_networks.0.cloud_servers.0.id", "ah_cloud_server.web", "id"),
					resource.TestCheckResourceAttrPair("data.ah_private_networks.test", "private_networks.0.cloud_servers.0.ip", "ah_private_network_connection.example", "ip_address"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
