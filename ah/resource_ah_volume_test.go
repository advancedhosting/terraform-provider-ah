package ah

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAHVolume_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_volume.test", "id"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "name"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "size"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "file_system"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "state"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAHVolume_CreateWithPlanSlug(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigCreateWithPlanSlug(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_volume.test", "id"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "name"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "size"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "file_system"),
				),
			},
		},
	})
}

func TestAccAHVolume_CreateWithoutFileSystem(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigCreateWithoutFileSystem(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_volume.test", "file_system", ""),
				),
			},
		},
	})
}

func TestAccAHVolume_CreateFromOrigin(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigFromOrigin(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_volume.test", "id"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "name"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "size"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "file_system"),
					resource.TestCheckResourceAttrSet("ah_volume.test", "origin_volume_id"),
				),
			},
		},
	})
}

func TestAccAHVolume_ChangeName(t *testing.T) {
	var beforeID, afterID string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeConfigChangeName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &afterID),
					resource.TestCheckResourceAttr("ah_volume.test", "name", "New Volume Name"),
					testAccCheckAHVolumeNoRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHVolume_IncreaseSize(t *testing.T) {
	var beforeID, afterID string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeConfigChangeSize(30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &afterID),
					resource.TestCheckResourceAttr("ah_volume.test", "size", "30"),
					testAccCheckAHVolumeNoRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHVolume_DowngradeSize(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigBasic(),
			},
			{
				Config:      testAccCheckAHVolumeConfigChangeSize(10),
				ExpectError: regexp.MustCompile("New size value must be greater than old value*"),
			},
		},
	})
}

func TestAccAHVolume_ChangeFileSystem(t *testing.T) {
	var beforeID, afterID string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHVolumeConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHVolumeConfigChangeFileSystem(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHVolumeExists("ah_volume.test", &afterID),
					resource.TestCheckResourceAttr("ah_volume.test", "file_system", "xfs"),
					testAccCheckAHVolumeRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func testAccCheckAHVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_volume" {
			continue
		}

		_, err := client.Volumes.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing volume (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHVolumeConfigBasic() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
	}`, VolumePlanID)
}

func testAccCheckAHVolumeConfigCreateWithSlug() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
	}`, VolumePlanName)
}

func testAccCheckAHVolumeConfigCreateWithPlanSlug() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		plan = "%s"
		file_system = "ext4"
		size = "20"
	}`, VolumePlanName)
}

func testAccCheckAHVolumeConfigCreateWithoutFileSystem() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = ""
		size = "20"
	}`, VolumePlanName)
}

func testAccCheckAHVolumeConfigFromOrigin() string {
	return fmt.Sprintf(`
	resource "ah_volume" "origin" {
		name = "Origin Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
	}

	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
		origin_volume_id = ah_volume.origin.id
	}`, VolumePlanName, VolumePlanName)
}

func testAccCheckAHVolumeConfigChangeName() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "New Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "20"
	}`, VolumePlanName)
}

func testAccCheckAHVolumeConfigChangeSize(newSize int) string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "ext4"
		size = "%d"
	}`, VolumePlanID, newSize)
}

func testAccCheckAHVolumeConfigChangeFileSystem() string {
	return fmt.Sprintf(`
	resource "ah_volume" "test" {
		name = "Volume Name"
		product = "%s"
		file_system = "xfs"
		size = "20"
	}`, VolumePlanName)
}

func testAccCheckAHVolumeExists(n string, volumeID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}

		*volumeID = rs.Primary.ID
		return nil
	}
}

func testAccCheckAHVolumeNoRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID != *afterID {
			t.Fatalf("Resource has been recreated, old ID: %s, new ID: %s", *beforeID, *afterID)
		}
		return nil
	}
}

func testAccCheckAHVolumeRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID == *afterID {
			t.Fatalf("Resource hasn't been recreated, ID: %s", *beforeID)
		}
		return nil
	}
}
