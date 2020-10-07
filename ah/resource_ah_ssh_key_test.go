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

func TestAccAHSSHKey_Basic(t *testing.T) {
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHSSHKeyConfigBasic(name, publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "id"),
					resource.TestCheckResourceAttr("ah_ssh_key.test", "name", name),
					resource.TestCheckResourceAttr("ah_ssh_key.test", "public_key", publicKey),
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "fingerprint"),
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "created_at"),
				),
			},
		},
	})
}

func TestAccAHSSHKey_EmptyName(t *testing.T) {
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHSSHKeyConfigEmptyName(publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "id"),
					resource.TestCheckResourceAttr("ah_ssh_key.test", "name", "test@ah-test.com"),
				),
			},
		},
	})
}

func TestAccAHSSHKey_UpdatePublicKey(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}
	newPublicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHSSHKeyConfigBasic(name, publicKey),
			},
			{
				Config: testAccCheckAHSSHKeyConfigBasic(name, newPublicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "id"),
					resource.TestCheckResourceAttr("ah_ssh_key.test", "public_key", newPublicKey),
				),
			},
		},
	})
}

func TestAccAHSSHKey_UpdateName(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	newName := fmt.Sprintf("test-%s", acctest.RandString(10))
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAHSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHSSHKeyConfigBasic(name, publicKey),
			},
			{
				Config: testAccCheckAHSSHKeyConfigBasic(newName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ah_ssh_key.test", "id"),
					resource.TestCheckResourceAttr("ah_ssh_key.test", "name", newName),
				),
			},
		},
	})
}

func testAccCheckAHSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_ssh_key" {
			continue
		}

		_, err := client.SSHKeys.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("Error removing volume (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHSSHKeyConfigBasic(name string, publicKey string) string {
	return fmt.Sprintf(`
	resource "ah_ssh_key" "test" {
	  name = "%s"
	  public_key = "%s"
	}`, name, publicKey)
}

func testAccCheckAHSSHKeyConfigEmptyName(publicKey string) string {
	return fmt.Sprintf(`
	resource "ah_ssh_key" "test" {
	  public_key = "%s"
	}`, publicKey)
}
