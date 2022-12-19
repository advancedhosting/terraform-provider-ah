package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAHIPAssignment_Basic(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHIPAssignmentConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "false"),
				),
			},
		},
	})
}

func TestAccAHIPAssignment_BasicByIP(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHIPAssignmentConfigBasicIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "false"),
				),
			},
		},
	})
}

func TestAccAHIPAssignment_Primary(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHIPAssignmentConfigPrimaryIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "true"),
				),
			},
		},
	})
}

func TestAccAHIPAssignment_UpdateToPrimary(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHIPAssignmentConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "false"),
				),
			},
			{
				Config: testAccCheckAHIPAssignmentConfigPrimaryIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "true"),
				),
			},
			{
				Config: testAccCheckAHIPAssignmentConfigDeleteAssignment(name),
			},
		},
	})
}

func testAccCheckAHIPAssignmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_ip_assignment" {
			continue
		}

		_, err := client.IPAddresses.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing ip assignment (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHIPAssignmentConfigBasic(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	  product = "start-xs"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.id
	}`, cloudServerName)
}

func testAccCheckAHIPAssignmentConfigBasicIP(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	  product = "start-xs"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.ip_address
	}`, cloudServerName)
}

func testAccCheckAHIPAssignmentConfigPrimaryIP(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	  product = "start-xs"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.id
	  primary = true
	}`, cloudServerName)
}

func testAccCheckAHIPAssignmentConfigDeleteAssignment(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	  product = "start-xs"
	}`, cloudServerName)
}
