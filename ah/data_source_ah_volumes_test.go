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
