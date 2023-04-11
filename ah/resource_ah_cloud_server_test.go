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
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "state", "running"),
					resource.TestCheckResourceAttrSet("ah_cloud_server.web", "plan"),
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

	datasourceConfig := fmt.Sprintf(`
	data "ah_cloud_images" "test" {
		filter {
			key = "slug"
			values = ["%s"]
		  }
	}`, ImageName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig + testAccCheckAHCloudServerConfigWithSSHKeys(name, publicKey, secondPublicKey),
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

	datasourceConfig := fmt.Sprintf(`
	data "ah_cloud_images" "test" {
		filter {
			key = "slug"
			values = ["%s"]
		  }
	}`, ImageName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig + testAccCheckAHCloudServerConfigWithoutPublicIP(name),
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
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigBasic(newName),
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
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "plan", VpsPlanID),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigUpgrade(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "plan", VpsUpgPlanID),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &afterID),
					testAccCheckAHResourceNoRecreated(t, &beforeID, &afterID),
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
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigCreateWithSlugs(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
					resource.TestCheckResourceAttr("ah_cloud_server.web", "plan", "start-xs"),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigUpgradeWithSlug(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "plan", VpsUpgPlanName),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &afterID),
					testAccCheckAHResourceNoRecreated(t, &beforeID, &afterID),
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
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_cloud_server.web", "name", name),
					testAccCheckAHCloudServerExists("ah_cloud_server.web", &beforeID),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHCloudServerConfigUpdateImageID(name),
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
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   plan = "%s"
	 }`, name, DatacenterID, VpsPlanID)
}

func testAccCheckAHCloudServerConfigCreateWithSlugs(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "%s"
	   plan = "%s"
	 }`, name, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerConfigCreateWithPlan(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "%s"
	   plan = "%s"
	 }`, name, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerConfigCreateWithAutoBackups(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "%s"
	   plan = "%s"
	  backups = true
	}`, name, DatacenterName, ImageName, VpsPlanName)
}

func testAccCheckAHCloudServerConfigCreateInPrivateCloud(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "%s"
	  private_cloud = true
	   node_id = "%s"
	  cluster_id = "%s"
	  vcpu = 1
	   ram = 64
	   disk = 10
	 }`, name, DatacenterName, ImageName, NodeID, ClusterID)
}

func testAccCheckAHCloudServerConfigUpgradeWithSlug(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "%s"
	   plan = "%s"
	 }`, name, DatacenterName, ImageName, VpsUpgPlanName)
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
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   plan = "%s"
	   ssh_keys = [ah_ssh_key.ssh_key1.id, ah_ssh_key.ssh_key2.fingerprint]
	 }`, ssh1PublicKey, ssh2PublicKey, name, DatacenterID, VpsPlanName)
}

func testAccCheckAHCloudServerConfigWithoutPublicIP(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   plan = "%s"
	   create_public_ip_address = false
	 }`, name, DatacenterID, VpsPlanName)
}

func testAccCheckAHCloudServerConfigUpgrade(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   plan = "%s"
	 }`, name, DatacenterID, VpsUpgPlanID)
}

func testAccCheckAHCloudServerConfigUpdateImageID(name string) string {
	return fmt.Sprintf(`
	 resource "ah_cloud_server" "web" {
	   name = "%s"
	   datacenter = "%s"
	   image = "8ed8bea7-69f0-40de-ab07-6a6b5a13581d"
	   plan = "%s"
	 }`, name, DatacenterID, VpsPlanName)
}

func testAccCheckAHCloudServerExists(n string, instanceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no instance ID is set")
		}

		*instanceID = rs.Primary.ID
		return nil
	}
}
