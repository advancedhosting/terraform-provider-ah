package ah

import (
	"context"
	"fmt"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAHCloudServer_Basic(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "state", "running"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "product"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "vcpu"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "ram"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "disk"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "datacenter"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "image"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "created_at"),

					resource.TestCheckResourceAttr("ah_cloud_server.web", "ips.#", "1"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "ips.0.assignment_id"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "ips.0.ip_address"),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "ips.0.primary", "true"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "ips.0.type"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "ips.0.reverse_dns"),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithSSHKey(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigWithSSHKey(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "ssh_keys.0", "232e3378-ff15-4f87-b63a-9ae6b34966a7"),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithoutPublicIP(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigWithoutPublicIP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("ah_cloud_server.web", "ips.0"),
				),
			},
		},
	})
}

func TestAccAHCloudServer_Rename(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	newName := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
				),
			},
			{
				Config: testAccCheckAHCloudServerConfigBasic(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", newName),
				),
			},
		},
	})
}

func TestAccAHCloudServer_Upgrade(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "df42a96b-b381-412c-a605-d66d7bf081af"),
				),
			},
			{
				Config: testAccCheckAHCloudServerConfigUpgrade(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "3ca84dd3-e439-46f4-8f47-f0fbb810896e"),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &afterID),
					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
				),
			},
		},
	})
}

func TestAccAHCloudServer_UpdateImage(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
				),
			},
			{
				Config: testAccCheckAHCloudServerConfigUpdateImageID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "image", "52ed921b-b5ca-4a5f-a3c9-69e283a126bf"),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &afterID),
					testAccCheckAHResourceRecreated(t, &beforeID, &afterID),
				),
			},
		},
	})
}

func testAccCheckAHCloudServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_cloud_server" {
			continue
		}

		_, err := client.Instances.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf(
				"Error waiting for instance (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHCloudServerConfigBasic(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	   product = "df42a96b-b381-412c-a605-d66d7bf081af"
	 }`, name)
}
func testAccCheckAHCloudServerConfigWithSSHKey(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	   product = "df42a96b-b381-412c-a605-d66d7bf081af"
	   ssh_keys = ["232e3378-ff15-4f87-b63a-9ae6b34966a7"]
	 }`, name)
}

func testAccCheckAHCloudServerConfigWithoutPublicIP(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	   product = "df42a96b-b381-412c-a605-d66d7bf081af"
	   create_public_ip_address = false
	 }`, name)
}

func testAccCheckAHCloudServerConfigUpgrade(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
	   product = "3ca84dd3-e439-46f4-8f47-f0fbb810896e"
	 }`, name)
}

func testAccCheckAHCloudServerConfigUpdateImageID(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "52ed921b-b5ca-4a5f-a3c9-69e283a126bf"
	   product = "df42a96b-b381-412c-a605-d66d7bf081af"
	 }`, name)
}

func testAccCheckAHCloudServerExists(n string, instanceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		*instanceID = rs.Primary.ID
		return nil
	}
}

func testAccCheckAHResourceNoRecreated(t *testing.T, beforeID, afterID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if beforeID != afterID {
			t.Fatalf("Resource has been recreated, old ID: %s, new ID: %s", beforeID, afterID)
		}
		return nil
	}
}

func testAccCheckAHResourceRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if beforeID == afterID {
			t.Fatalf("Resource hasn't been recreated, ID: %s", *beforeID)
		}
		return nil
	}
}
