package ah

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHCloudServerProducts() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return errors.New("use ah_cloud_server_plans resource instead")
		},
	}
}
