package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAHPrivateNetworkConnection_Basic(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHPrivateNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_private_network_connection.example", "id"),
					resource.TestCheckResourceAttrPair("ah_private_network_connection.example", "cloud_server_id", "ah_cloud_server.web", "id"),
					resource.TestCheckResourceAttrPair("ah_private_network_connection.example", "private_network_id", "ah_private_network.test", "id"),
					resource.TestCheckResourceAttrSet("ah_private_network_connection.example", "ip_address"),
				),
			},
		},
	})
}

func TestAccAHPrivateNetworkConnection_UpdateIP(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHPrivateNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigUpdateIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_private_network_connection.example", "ip_address", "10.0.0.2"),
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &afterID),
					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
				),
			},
		},
	})
}

func TestAccAHPrivateNetworkConnection_ChangeCloudServer(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHPrivateNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigUpdateCloudServer(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("ah_private_network_connection.example", "cloud_server_id", "ah_cloud_server.web.1", "id"),
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &afterID),
					testAccCheckAHResourceRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func TestAccAHPrivateNetworkConnection_ChangePrivateNetwork(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHPrivateNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPrivateNetworkConnectionConfigUpdatePrivateNetwork(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("ah_private_network_connection.example", "private_network_id", "ah_private_network.test.1", "id"),
					testAccCheckAHPrivateNetworkConnectionExists("ah_private_network_connection.example", &afterID),
					testAccCheckAHResourceRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func testAccCheckAHPrivateNetworkConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_private_network_connection" {
			continue
		}

		_, err := client.IPAddresses.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing private network connection (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHPrivateNetworkConnectionConfigBasic(cloudServerName string) string {
	return fmt.Sprintf(`
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "df42a96b-b381-412c-a605-d66d7bf081af"
	}
	
	resource "ah_private_network_connection" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  private_network_id = ah_private_network.test.id
	}`, cloudServerName)
}

func testAccCheckAHPrivateNetworkConnectionConfigUpdateIP(cloudServerName string) string {
	return fmt.Sprintf(`
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "df42a96b-b381-412c-a605-d66d7bf081af"
	}
	
	resource "ah_private_network_connection" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  private_network_id = ah_private_network.test.id
	  ip_address = "10.0.0.2"
	}`, cloudServerName)
}

func testAccCheckAHPrivateNetworkConnectionConfigUpdateCloudServer(cloudServerName string) string {
	return fmt.Sprintf(`
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}
	 
	resource "ah_cloud_server" "web" {
	  count = 2
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "df42a96b-b381-412c-a605-d66d7bf081af"
	}
	
	resource "ah_private_network_connection" "example" {
	  cloud_server_id = ah_cloud_server.web.1.id
	  private_network_id = ah_private_network.test.id
	}`, cloudServerName)
}

func testAccCheckAHPrivateNetworkConnectionConfigUpdatePrivateNetwork(cloudServerName string) string {
	return fmt.Sprintf(`
	resource "ah_private_network" "test" {
	  count = 2
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}
	 
	resource "ah_cloud_server" "web" {
	  name = "%s"
	  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	  product = "df42a96b-b381-412c-a605-d66d7bf081af"
	}
	
	resource "ah_private_network_connection" "example" {
	  cloud_server_id = ah_cloud_server.web.id
	  private_network_id = ah_private_network.test.1.id
	}`, cloudServerName)
}

func testAccCheckAHPrivateNetworkConnectionExists(n string, privateNetworkID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance private network ID is set")
		}

		*privateNetworkID = rs.Primary.ID
		return nil
	}
}
