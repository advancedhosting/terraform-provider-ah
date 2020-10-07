package ah

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAHIP_BasicPublicIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPublicIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "ip_address"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "created_at"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "reverse_dns"),
				),
			},
		},
	})
}
func TestAccAHIP_BasicPublicIPWithoutDatacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAHPublicIPConfigWithoutdDatacenter(),
				ExpectError: regexp.MustCompile(`.*Datacenter is required for public ip.*`),
			},
		},
	})
}

func TestAccAHIP_PublicIPWithReserveDNS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHPublicIPConfigBasicWithReserveDNS(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "ip_address"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "created_at"),
					resource.TestCheckResourceAttr("ah_ip.test", "reverse_dns", "ip-185-189-69-16.ah-server22.com"),
				),
			},
		},
	})
}

func TestAccAHIP_BasicAnycastIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHAnycastIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ip.test", "id"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "ip_address"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "created_at"),
					resource.TestCheckResourceAttrSet("ah_ip.test", "reverse_dns"),
					resource.TestCheckNoResourceAttr("ah_ip.test", "datacenter"),
				),
			},
		},
	})
}

func TestAccAHIP_UpdateReverseDNS(t *testing.T) {
	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHAnycastIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPublicIPConfigBasicWithReserveDNS(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &afterID),
					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
					resource.TestCheckResourceAttr("ah_ip.test", "reverse_dns", "ip-185-189-69-16.ah-server22.com"),
				),
			},
		},
	})
}

func TestAccAHIP_UpdateDatacenter(t *testing.T) {
	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHAnycastIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHPublicIPConfigNewDatacenter(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &afterID),
					testAccCheckAHResourceRecreated(t, &beforeID, &afterID),
					resource.TestCheckResourceAttr("ah_ip.test", "datacenter", "1b1ae192-d44e-451b-8d39-a8670c58e97d"),
				),
			},
		},
	})
}

func TestAccAHIP_UpdateType(t *testing.T) {
	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHAnycastIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHAnycastIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHIPExists("ah_ip.test", &afterID),
					testAccCheckAHResourceRecreated(t, &beforeID, &afterID),
					resource.TestCheckResourceAttr("ah_ip.test", "type", "anycast"),
				),
			},
		},
	})
}

func testAccCheckAHIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_ip" {
			continue
		}

		_, err := client.IPAddresses.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing ip (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHPublicIPConfigBasic() string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	 }`)
}

func testAccCheckAHPublicIPConfigNewDatacenter() string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "1b1ae192-d44e-451b-8d39-a8670c58e97d"
	 }`)
}

func testAccCheckAHPublicIPConfigWithoutdDatacenter() string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	 }`)
}

func testAccCheckAHPublicIPConfigBasicWithReserveDNS() string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "public"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   reverse_dns = "ip-185-189-69-16.ah-server22.com"
	 }`)
}

func testAccCheckAHAnycastIPConfigBasic() string {
	return fmt.Sprintf(`
	 resource "ah_ip" "test" {
	   type = "anycast"
	 }`)
}

func testAccCheckAHIPExists(n string, ipID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ip ID is set")
		}

		*ipID = rs.Primary.ID
		return nil
	}
}
