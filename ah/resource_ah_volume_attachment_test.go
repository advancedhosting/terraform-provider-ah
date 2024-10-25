package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAHVolumeAttachment_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_volume_attachment.test", "id"),
					resource.TestCheckResourceAttrPair("ah_volume_attachment.test", "cloud_server_id", "ah_cloud_server.web", "id"),
					resource.TestCheckResourceAttrPair("ah_volume_attachment.test", "volume_id", "ah_volume.test", "id"),
					resource.TestCheckResourceAttrSet("ah_volume_attachment.test", "state"),
				),
			},
		},
	})
}

func TestAccAHVolume_IncreaseSizeAttachedVolume(t *testing.T) {
	var beforeID, afterID string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &beforeID),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigBasic(30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &afterID),
					resource.TestCheckResourceAttr("ah_volume.test", "size", "30"),
					testAccCheckAHVolumeNoRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHVolumeAttachment_ChangeCloudServer(t *testing.T) {
	var beforeID, afterID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &beforeID),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigChangeCloudServer(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &afterID),
					testAccCheckAHVolumeAttachmentRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHVolumeAttachment_ChangeVolume(t *testing.T) {
	var beforeID, afterID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &beforeID),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHVolumeAttachmentConfigChangeVolume(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &afterID),
					testAccCheckAHVolumeAttachmentRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func testAccCheckAHVolumeAttachmentExists(n string, volumeAttachmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume attachment ID is set")
		}

		*volumeAttachmentID = rs.Primary.ID
		return nil
	}
}

func testAccCheckAHVolumeAttachmentRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID == *afterID {
			t.Fatalf("Resource hasn't been recreated, ID: %s", *beforeID)
		}
		return nil
	}
}

func testAccCheckAHVolumeAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_volume_attachment" {
			continue
		}

		cloudServerID := rs.Primary.Attributes["cloud_server_id"]
		volumeID := rs.Primary.Attributes["volume_id"]
		instance, err := client.Instances.Get(context.Background(), cloudServerID)

		if err == ah.ErrResourceNotFound {
			return nil
		}

		if err != nil {
			return fmt.Errorf("Error removing volume %s from cloud server %s: %s", volumeID, cloudServerID, err)
		}

		for _, volume := range instance.Volumes {
			if volume.ID == volumeID {
				return fmt.Errorf("Volume %s has not been detached from cloud server %s", volumeID, cloudServerID)
			}
		}
	}

	return nil
}

func testAccCheckAHVolumeAttachmentConfigBasic(volumeSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "%d"
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
	}`, VolumePlanID, volumeSize, DatacenterID, VpsPlanID)
}

func testAccCheckAHVolumeAttachmentConfigChangeCloudServer(volumeSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "%d"
	}
	 
	resource "ah_cloud_server" "web" {
	  count = 2
	  name = "test"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}
	
	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.1.id
	  volume_id = ah_volume.test.id
	}`, VolumePlanID, volumeSize, DatacenterID, VpsPlanID)
}

func testAccCheckAHVolumeAttachmentConfigChangeVolume(volumeSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		count = 2
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "%d"
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "b63d5134-4a1b-4329-ab23-e6597daa53a5"
	  product = "%s"
	}
	
	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  volume_id = ah_volume.test.1.id
	}`, VolumePlanID, volumeSize, VpsPlanID)
}
