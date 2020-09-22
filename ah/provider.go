package ah

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{

		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AH_ACCESS_TOKEN", nil),
				Description: "The API token to access the AH cloud.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AH_API_URL", nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ah_cloud_servers":    dataSourceAHCloudServers(),
			"ah_ips":              dataSourceAHIPs(),
			"ah_private_networks": dataSourceAHPrivateNetworks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ah_cloud_server":               resourceAHCloudServer(),
			"ah_ip":                         resourceAHIP(),
			"ah_ip_assignment":              resourceAHIPAssignment(),
			"ah_private_network":            resourceAHPrivateNetwork(),
			"ah_private_network_connection": resourceAHPrivateNetworkConnection(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {

	config := Config{
		Token:       d.Get("access_token").(string),
		APIEndpoint: d.Get("endpoint").(string),
	}

	return config.Client()
}
