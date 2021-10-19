package ah

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAHSSHKeys_Basic(t *testing.T) {
	publicKey, _, err := acctest.RandSSHKeyPair("test@ah-test.com")
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	if err != nil {
		t.Fatalf("RandSSHKeyPair error: %s", err)
	}

	resourcesConfig := testAccCheckAHSSHKeyConfigBasic(name, publicKey)

	datasourceConfig := `
	data "ah_ssh_keys" "test" {
		filter {
			key = "name"
			values = [ah_ssh_key.test.name]
		}
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
