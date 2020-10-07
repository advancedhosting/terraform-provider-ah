package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAHVolumeAttachment_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeAttachmentConfigBasic(20),
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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeAttachmentConfigBasic(30),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeAttachmentConfigChangeCloudServer(20),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeAttachmentConfigBasic(20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeAttachmentExists("ah_volume_attachment.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeAttachmentConfigChangeVolume(20),
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
		product = "03bebb65-22d8-43c6-819b-5b85b5e49c82"
		file_system = "ext4"
		size = "%d"
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
	}`, volumeSize)
}

func testAccCheckAHVolumeAttachmentConfigChangeCloudServer(volumeSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "03bebb65-22d8-43c6-819b-5b85b5e49c82"
		file_system = "ext4"
		size = "%d"
	}
	 
	resource "ah_cloud_server" "web" {
	  count = 2
	  name = "test"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "1a4cdeb2-6ca4-4745-819e-ac2ea99dc0cc"
	}
	
	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.1.id
	  volume_id = ah_volume.test.id
	}`, volumeSize)
}

func testAccCheckAHVolumeAttachmentConfigChangeVolume(volumeSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		count = 2
		name = "Volume Name"
		product = "03bebb65-22d8-43c6-819b-5b85b5e49c82"
		file_system = "ext4"
		size = "%d"
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "1a4cdeb2-6ca4-4745-819e-ac2ea99dc0cc"
	}
	
	resource "ah_volume_attachment" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  volume_id = ah_volume.test.1.id
	}`, volumeSize)
}
