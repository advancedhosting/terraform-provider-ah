package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHCloudServers_Basic(t *testing.T) {

	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resourcesConfig := fmt.Sprintf(`
	resource "ah_cloud_server" "web" {
	  count = 3
	  name = "%s_${count.index}"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	  product = "start-xs"
	}`, name)

	datasourceConfig := fmt.Sprintf(`
	data "ah_cloud_servers" "test" {
		filter {
			key = "name"
			values = ["%s_1"]
		}
		sort {
			key = "created_at"
			direction = "desc"
		}
	}`, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_cloud_servers.test", "cloud_servers.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_servers.test", "cloud_servers.0.id", "ah_cloud_server.web.1", "id"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_servers.test", "cloud_servers.0.name", "ah_cloud_server.web.1", "name"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.datacenter"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.product"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.state"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.vcpu"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ram"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.disk"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.image"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.backups"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.use_password"),
					resource.TestCheckResourceAttr("data.ah_cloud_servers.test", "cloud_servers.0.ips.#", "1"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ips.0.assignment_id"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ips.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ips.0.primary"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ips.0.type"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_servers.test", "cloud_servers.0.ips.0.reverse_dns"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
