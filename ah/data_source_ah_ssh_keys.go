package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHSSHKeys() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "fingerprint"}
	allowedSortingKeys := []string{"id", "name", "fingerprint", "created_at"}
	return &schema.Resource{
		Read: dataSourceAHSSHKeysRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"ssh_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAHSSHKeysListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		sorting := &ah.Sorting{
			Key:   m["key"].(string),
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHSSHKeysListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		filter := &ah.InFilter{
			Keys:   []string{m["key"].(string)},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHSSHKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHSSHKeysListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHSSHKeysListSorting(v.(*schema.Set))
	}

	sshKeys, err := allSSHKeysInfo(client, options)
	if err != nil {
		return err
	}

	if err = dataSourceAHSSHKeysSchema(d, meta, sshKeys); err != nil {
		return err
	}
	return nil
}

func dataSourceAHSSHKeysSchema(d *schema.ResourceData, meta interface{}, sshKeys []ah.SSHKey) error {
	allSSHKeys := make([]map[string]interface{}, len(sshKeys))
	var ids string
	for i, sshKey := range sshKeys {
		sshKeyInfo := map[string]interface{}{
			"id":          sshKey.ID,
			"name":        sshKey.Name,
			"public_key":  sshKey.PublicKey,
			"fingerprint": sshKey.Fingerprint,
			"created_at":  sshKey.CreatedAt,
		}
		ids += sshKey.ID
		allSSHKeys[i] = sshKeyInfo
	}
	if err := d.Set("ssh_keys", allSSHKeys); err != nil {
		return fmt.Errorf("unable to set ssh_keys attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}

func allSSHKeysInfo(client *ah.APIClient, options *ah.ListOptions) ([]ah.SSHKey, error) {
	meta := &ah.ListMetaOptions{
		Page: 1,
	}

	options.Meta = meta

	var allSSHKeys []ah.SSHKey

	for {
		sshKeys, meta, err := client.SSHKeys.List(context.Background(), options)

		if err != nil {
			return nil, fmt.Errorf("Error list ssh keys: %s", err)
		}

		allSSHKeys = append(allSSHKeys, sshKeys...)
		if meta.IsLastPage() {
			break
		}

		options.Meta.Page++
	}

	return allSSHKeys, nil
}
