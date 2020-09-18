package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAHSSHKeys_Basic(t *testing.T) {
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}

	resourcesConfig := testAccCheckAHSSHKeyConfigBasic(name, publicKey)

	// TODO add filter after WCS-3609
	datasourceConfig := `
	data "ah_ssh_keys" "test" {}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ah_ssh_keys.test", "ssh_keys.0.id"),
					resource.TestCheckResourceAttrSet("data.ah_ssh_keys.test", "ssh_keys.0.name"),
					resource.TestCheckResourceAttrSet("data.ah_ssh_keys.test", "ssh_keys.0.public_key"),
					resource.TestCheckResourceAttrSet("data.ah_ssh_keys.test", "ssh_keys.0.fingerprint"),
					resource.TestCheckResourceAttrSet("data.ah_ssh_keys.test", "ssh_keys.0.created_at"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
