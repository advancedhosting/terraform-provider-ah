package ah

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var NodePoolSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
	},
	"nodes_count": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"created_at": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"auto_scale": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"min_count": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"max_count": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"public_properties": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
	},
	"private_properties": {
		Type:     schema.TypeMap,
		Optional: true,
	},
}
