package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAHCloudServerSnapshot_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "id"),
					resource.TestCheckResourceAttrPair("ah_cloud_server_snapshot.test", "cloud_server_id", "ah_cloud_server.web", "id"),
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "state"),
					resource.TestCheckResourceAttr("ah_cloud_server_snapshot.test", "name", "example-snapshot-1"),
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "size"),
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "type"),
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAHCloudServerSnapshot_EmptyName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerSnapshotConfigEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_cloud_server_snapshot.test", "name"),
				),
			},
		},
	})
}

func TestAccAHCloudServerSnapshot_UpdateName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var beforeID, afterID string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHCloudServerSnapshotExists("ah_cloud_server_snapshot.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHCloudServerSnapshotConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHCloudServerSnapshotExists("ah_cloud_server_snapshot.test", &afterID),
					resource.TestCheckResourceAttr("ah_cloud_server_snapshot.test", "name", "New Name"),
					testAccCheckAHCloudServerSnapshotNoRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHCloudServerSnapshot_UpdateCloudServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerSnapshotConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHCloudServerSnapshotExists("ah_cloud_server_snapshot.test", &beforeID),
					resource.TestCheckResourceAttrPair("ah_cloud_server_snapshot.test", "cloud_server_id", "ah_cloud_server.web", "id"),
				),
			},
			{
				Config: testAccCheckAHCloudServerSnapshotConfigUpdateCloudServer(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHCloudServerSnapshotExists("ah_cloud_server_snapshot.test", &afterID),
					resource.TestCheckResourceAttrPair("ah_cloud_server_snapshot.test", "cloud_server_id", "ah_cloud_server.web.1", "id"),
					testAccCheckAHCloudServerSnapshotRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func testAccCheckAHCloudServerSnapshotExists(n string, BackupID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume attachment ID is set")
		}

		*BackupID = rs.Primary.ID
		return nil
	}
}

func testAccCheckAHCloudServerSnapshotNoRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID != *afterID {
			t.Fatalf("Resource has been recreated, expected %s, got %s", *beforeID, *afterID)
		}
		return nil
	}
}

func testAccCheckAHCloudServerSnapshotRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID == *afterID {
			t.Fatalf("Resource hasn't been recreated")
		}
		return nil
	}
}

func testAccCheckAHCloudServerSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_cloud_server_snapshot" {
			continue
		}

		_, err := client.Backups.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing backup %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckAHCloudServerSnapshotConfigBasic() string {
	return fmt.Sprintf(`
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "%s"
	  image = "%s"
	  product = "%s"
	}

	resource "ah_cloud_server_snapshot" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  name = "example-snapshot-1"
	}`, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerSnapshotConfigEmptyName() string {
	return fmt.Sprintf(`
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "%s"
	  image = "%s"
	  product = "%s"
	}
	
	resource "ah_cloud_server_snapshot" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	}`, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerSnapshotConfigUpdateName() string {
	return fmt.Sprintf(`
	resource "ah_cloud_server" "web" {
	  name = "test"
	  datacenter = "%s"
	  image = "%s"
	  product = "%s"
	}
	
	resource "ah_cloud_server_snapshot" "test" {
	  cloud_server_id = ah_cloud_server.web.id
	  name = "New Name"
	}`, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerSnapshotConfigUpdateCloudServer() string {
	return fmt.Sprintf(`
	resource "ah_cloud_server" "web" {
	  count = 2
	  name = "test"
	  datacenter = "%s"
	  image = "%s"
	  product = "%s"
	}
	
	resource "ah_cloud_server_snapshot" "test" {
	  cloud_server_id = ah_cloud_server.web.1.id
	  name = "example-snapshot-1"
	}`, DatacenterName, ImageName, VpsPlanName)
}
