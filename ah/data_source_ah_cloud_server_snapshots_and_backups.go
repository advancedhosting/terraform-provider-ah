package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAHCloudServerSnapshotsAndBackups() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "cloud_server_id", "cloud_server_name", "state", "size", "type"}
	allowedSortingKeys := []string{"id", "name", "cloud_server_id", "cloud_server_name", "state", "size", "type", "created_at"}
	return &schema.Resource{
		Read: dataSourceAHCloudServerSnapshotsAndBackupsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFilterSchema(allowedFilterKeys),
			"sort":   dataSourceSortingSchema(allowedSortingKeys),
			"snapshots_and_backups": {
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
						"cloud_server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
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

func buildAHCloudServerSnapshotsAndBackupsListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		key := m["key"].(string)

		switch key {
		case "cloud_server_id":
			key = "instance_id"
		case "state":
			key = "status"
		case "cloud_server_name":
			key = "instance_name"
		}

		sorting := &ah.Sorting{
			Key:   key,
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func buildAHCloudServerSnapshotsAndBackupsListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		key := m["key"].(string)

		switch key {
		case "cloud_server_id":
			key = "instance_id"
		case "state":
			key = "status"
		case "cloud_server_name":
			key = "instance_name"
		}

		filter := &ah.InFilter{
			Keys:   []string{key},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func dataSourceAHCloudServerSnapshotsAndBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{}

	if v, ok := d.GetOk("filter"); ok {
		options.Filters = buildAHCloudServerSnapshotsAndBackupsListFilter(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sortings = buildAHCloudServerSnapshotsAndBackupsListSorting(v.(*schema.Set))
	}

	instancesBackups, err := client.Backups.List(context.Background(), options)
	if err != nil {
		return err
	}

	if err = dataSourceAHCloudServerSnapshotsAndBackupsSchema(d, meta, instancesBackups); err != nil {
		return err
	}
	return nil
}

func dataSourceAHCloudServerSnapshotsAndBackupsSchema(d *schema.ResourceData, meta interface{}, instancesBackups []ah.InstanceBackups) error {
	var allBackups []map[string]interface{}
	var ids string
	for _, instanceBackup := range instancesBackups {
		for _, backup := range instanceBackup.Backups {
			backupInfo := map[string]interface{}{
				"id":                   backup.ID,
				"name":                 backup.Note,
				"cloud_server_id":      backup.InstanceID,
				"cloud_server_name":    instanceBackup.InstanceName,
				"cloud_server_deleted": instanceBackup.InstanceRemoved,
				"state":                backup.Status,
				"size":                 backup.Size,
				"type":                 backup.Type,
				"created_at":           backup.CreatedAt,
			}
			allBackups = append(allBackups, backupInfo)
			ids += backup.ID
		}
	}
	if err := d.Set("snapshots_and_backups", allBackups); err != nil {
		return fmt.Errorf("unable to set snapshots_and_backups attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}
