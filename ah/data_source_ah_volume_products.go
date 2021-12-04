package ah

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHVolumeProducts() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return errors.New("use ah_volume_plans resource instead")
		},
	}
}
