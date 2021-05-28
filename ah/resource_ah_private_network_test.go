package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAHPrivateNetwork_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHPrivateNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_private_network.test", "id"),
					resource.TestCheckResourceAttr("ah_private_network.test", "ip_range", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ah_private_network.test", "name", "Test Private Network"),
					resource.TestCheckResourceAttrSet("ah_private_network.test", "state"),
					resource.TestCheckResourceAttrSet("ah_private_network.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAHPrivateNetwork_UpdateName(t *testing.T) {
	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHPrivateNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPrivateNetworkConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHPrivateNetworkExists("ah_private_network.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPrivateNetworkConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHPrivateNetworkExists("ah_private_network.test", &afterID),
					resource.TestCheckResourceAttr("ah_private_network.test", "name", "New Private Network"),
					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
				),
			},
		},
	})
}

func testAccCheckAHPrivateNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_private_network" {
			continue
		}

		_, err := client.PrivateNetworks.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing private network (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHPrivateNetworkConfigBasic() string {
	return `
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "Test Private Network"
	}`
}

func testAccCheckAHPrivateNetworkConfigUpdateName() string {
	return `
	resource "ah_private_network" "test" {
	  ip_range = "10.0.0.0/24"
	  name = "New Private Network"
	}`
}

func testAccCheckAHPrivateNetworkExists(n string, privateNetworkID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No private network ID is set")
		}

		*privateNetworkID = rs.Primary.ID
		return nil
	}
}
