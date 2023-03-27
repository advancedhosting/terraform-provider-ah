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
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigBasic(name),
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
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigBasicIP(name),
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
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigPrimaryIP(name),
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
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "false"),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigPrimaryIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip_assignment.example", "id"),
					resource.TestCheckResourceAttr("ah_ip_assignment.example", "primary", "true"),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHIPAssignmentConfigDeleteAssignment(name),
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
	   datacenter = "%s"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.id
	}`, DatacenterID, cloudServerName, DatacenterID, VpsPlanID)
}

func testAccCheckAHIPAssignmentConfigBasicIP(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "%s"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.ip_address
	}`, DatacenterID, cloudServerName, DatacenterID, VpsPlanID)
}

func testAccCheckAHIPAssignmentConfigPrimaryIP(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "%s"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}
	
	resource "ah_ip_assignment" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  ip_address = ah_ip.test.id
	  primary = true
	}`, DatacenterID, cloudServerName, DatacenterID, VpsPlanID)
}

func testAccCheckAHIPAssignmentConfigDeleteAssignment(cloudServerName string) string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "%s"
	 }
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "%s"
	  image = "${data.ah_cloud_images.test.images.0.id}"
	  product = "%s"
	}`, DatacenterID, cloudServerName, DatacenterID, VpsPlanID)
}

func datasourceConfigBasic() string {
	return fmt.Sprintf(`
	data "ah_cloud_images" "test" {
		filter {
			key = "slug"
			values = ["%s"]
		  }
	}`, ImageName)
}
