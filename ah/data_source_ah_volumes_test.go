package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHVolumes_Basic(t *testing.T) {

	resourcesConfig := testAccCheckAHVolumeConfigBasic()

	datasourceConfig := `
	data "ah_volumes" "test" {
	  filter {
		key = "id"
		values = [ah_volume.test.id]
	  }
	  sort {
		key = "created_at"
		direction = "desc"
	  }
	}`

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
					resource.TestCheckResourceAttr("data.ah_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_volumes.test", "volumes.0.id", "ah_volume.test", "id"),
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.state"),
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.size"),
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.file_system"),
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.created_at"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}

func TestAccDataSourceAHVolumes_FilterByCloudServerID(t *testing.T) {

	resourcesConfig := fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
	}

	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}

	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  volume_id = ah_volume.test.id
	}
	`, "381347560", DatacenterID, VpsPlanName)

	datasourceConfig := `
	data "ah_volumes" "test" {
		filter {
		  key = "cloud_server_id"
		  values = [ah_cloud_server.web.id]
		}
		sort {
		  key = "created_at"
		  direction = "desc"
		}
	  }
	`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + datasourceConfig + resourcesConfig,
			},
			{
				Config: datasourceConfigBasic() + resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_volumes.test", "volumes.0.cloud_server_id", "ah_cloud_server.web", "id"),
				),
			},
			{
				Config: datasourceConfigBasic() + resourcesConfig,
			},
		},
	})
}
