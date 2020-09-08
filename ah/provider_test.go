package ah

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccProviders map[string]terraform.ResourceProvider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ah": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("AH_ACCESS_TOKEN") == "" {
		t.Fatal("AH_ACCESS_TOKEN must be set for acceptance tests")
	}
}
