package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
					resource.TestCheckResourceAttrSet("data.ah_volumes.test", "volumes.0.product"),
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

	resourcesConfig := `
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "03bebb65-22d8-43c6-819b-5b85b5e49c82"
		file_system = "ext4"
		size = 10
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "1a4cdeb2-6ca4-4745-819e-ac2ea99dc0cc"
	}
	
	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  volume_id = ah_volume.test.id
	}
	`

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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_volumes.test", "volumes.0.cloud_server_id", "ah_cloud_server.web", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
