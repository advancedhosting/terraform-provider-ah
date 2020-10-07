package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAHImages_Basic(t *testing.T) {

	datasourceConfig := `
	data "ah_cloud_images" "test" {
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.distribution"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.version"),
					resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.architecture"),
					//resource.TestCheckResourceAttrSet("data.ah_cloud_images.test", "images.0.slug"),
				),
			},
		},
	})
}
