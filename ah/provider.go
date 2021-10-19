package ah

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
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
			"ah_cloud_servers":                     dataSourceAHCloudServers(),
			"ah_ips":                               dataSourceAHIPs(),
			"ah_private_networks":                  dataSourceAHPrivateNetworks(),
			"ah_volumes":                           dataSourceAHVolumes(),
			"ah_cloud_server_snapshot_and_backups": dataSourceAHCloudServerSnapshotsAndBackups(),
			"ah_ssh_keys":                          dataSourceAHSSHKeys(),
			//"ah_volume_products":                   dataSourceAHVolumeProducts(),
			"ah_datacenters":  dataSourceAHDatacenters(),
			"ah_cloud_images": dataSourceAHImages(),
			//"ah_cloud_server_products":             dataSourceAHCloudServerProducts(),
			"ah_cloud_server_plans": dataSourceAHCloudServerPlans(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ah_cloud_server":               resourceAHCloudServer(),
			"ah_ip":                         resourceAHIP(),
			"ah_ip_assignment":              resourceAHIPAssignment(),
			"ah_private_network":            resourceAHPrivateNetwork(),
			"ah_private_network_connection": resourceAHPrivateNetworkConnection(),
			"ah_volume":                     resourceAHVolume(),
			"ah_volume_attachment":          resourceAHVolumeAttachment(),
			"ah_cloud_server_snapshot":      resourceAHCloudServerSnapshot(),
			"ah_ssh_key":                    resourceAHSSHKey(),
			"ah_load_balancer":              resourceAHLoadBalancer(),
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
