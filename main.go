package main

import (
	"github.com/advancedhosting/terraform-provider-ah/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ah.Provider})
}
