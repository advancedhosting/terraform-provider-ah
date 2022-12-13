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

func TestAccAHCloudServer_Basic(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
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

func TestAccAHCloudServer_CreateWithSlugs(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigCreateWithSlugs(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithPlan(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigCreateWithPlan(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithAutoBackups(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigCreateWithAutoBackups(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithSSHKey(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	publicKey, _, _ := acctest.RandSSHKeyPair("test@ah-test.com")
	secondPublicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigWithSSHKeys(name, publicKey, secondPublicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "ssh_keys.#", "2"),
					resource.TestCheckResourceAttrPair("ah_cloud_server.web", "ssh_keys.0", "ah_ssh_key.ssh_key1", "id"),
					resource.TestCheckResourceAttrPair("ah_cloud_server.web", "ssh_keys.1", "ah_ssh_key.ssh_key2", "fingerprint"),
				),
			},
		},
	})
}

func TestAccAHCloudServer_CreateWithoutPublicIP(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
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

func TestAccAHCloudServer_CreateInPrivateCloud(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigCreateInPrivateCloud(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "vcpu", "1"),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "ram", "64"),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "disk", "10"),
				),
			},
		},
	})
}

func TestAccAHCloudServer_Rename(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	newName := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "start-xs"),
				),
			},
			{
				Config: testAccCheckAHCloudServerConfigUpgrade(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "start-xs"),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &afterID),
					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
				),
			},
		},
	})
}

func TestAccAHCloudServer_UpgradeWithSlug(t *testing.T) {
	var beforeID, afterID string
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHCloudServerConfigCreateWithSlugs(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "start-xs"),
				),
			},
			{
				Config: testAccCheckAHCloudServerConfigUpgradeWithSlug(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "product", "start-m"),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
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
					resource.TestCheckResourceAttr("ah_cloud_server.web", "image", "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"),
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
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   product = "start-xs"
	 }`, name)
}

func testAccCheckAHCloudServerConfigCreateWithSlugs(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "ams1"
	   image = "ubuntu-20_04-x64"
	   product = "start-xs"
	 }`, name)
}

func testAccCheckAHCloudServerConfigCreateWithPlan(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "ams1"
	   image = "ubuntu-20_04-x64"
	   plan = "start-xs"
	 }`, name)
}

func testAccCheckAHCloudServerConfigCreateWithAutoBackups(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "ams1"
	   image = "ubuntu-20_04-x64"
	   product = "start-xs"
	   backups = true
	 }`, name)
}

func testAccCheckAHCloudServerConfigCreateInPrivateCloud(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "ams1"
	   image = "ubuntu-20_04-x64"
       private_cloud = true
	   node_id = "2486b2f8-f7a6-4207-979b-9b94d93c174e"
       cluster_id = "6770e666-7a7b-4e9f-816d-cfc98a52c84d"
       vcpu = 1
	   ram = 64
	   disk = 10
	 }`, name)
}

func testAccCheckAHCloudServerConfigUpgradeWithSlug(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "ams1"
	   image = "ubuntu-20_04-x64"
	   product = "start-m"
	 }`, name)
}

func testAccCheckAHCloudServerConfigWithSSHKeys(name string, ssh1PublicKey, ssh2PublicKey string) string {
	return fmt.Sprintf(`
	 resource "ah_ssh_key" "ssh_key1" {
	   name = "test"
	   public_key = "%s"
	 }
	 resource "ah_ssh_key" "ssh_key2" {
	   name = "test"
	   public_key = "%s"
	 }
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   product = "start-xs"
	   ssh_keys = [ah_ssh_key.ssh_key1.id, ah_ssh_key.ssh_key2.fingerprint]
	 }`, ssh1PublicKey, ssh2PublicKey, name)
}

func testAccCheckAHCloudServerConfigWithoutPublicIP(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   product = "start-xs"
	   create_public_ip_address = false
	 }`, name)
}

func testAccCheckAHCloudServerConfigUpgrade(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   product = "start-xs"
	 }`, name)
}

func testAccCheckAHCloudServerConfigUpdateImageID(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   product = "start-xs"
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
