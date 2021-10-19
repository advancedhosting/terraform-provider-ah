package ah

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAHCloudServerSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHCloudServerSnapshotCreate,
		Read:   resourceAHCloudServerSnapshotRead,
		Update: resourceAHCloudServerSnapshotUpdate,
		Delete: resourceAHCloudServerSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"cloud_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_server_name": {
				Type:     schema.TypeString,
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
	}
}

func resourceAHCloudServerSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	instanceID := d.Get("cloud_server_id").(string)

	var note string
	if attr, ok := d.GetOk("name"); ok {
		note = attr.(string)
	} else {
		note = time.Now().Format("2006-01-02 at 15:04:05")
	}

	action, err := client.Instances.CreateBackup(context.Background(), instanceID, note)

	if err != nil {
		return fmt.Errorf("Error creating backup: %s", err)
	}

	if err := waitForBackupReady(instanceID, action.ID, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for backup to become ready: %v", err)
	}

	action, err = client.Instances.ActionInfo(context.Background(), instanceID, action.ID)
	if err != nil {
		return fmt.Errorf("Error getting backup info: %s", err)
	}

	d.SetId(action.ResultParams.SnapshotID)
	return resourceAHCloudServerSnapshotRead(d, meta)

}

func resourceAHCloudServerSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	backup, instanceName, err := snapshotInfo(d, meta)
	if err != nil {
		return err
	}

	d.Set("cloud_server_id", backup.InstanceID)
	d.Set("cloud_server_name", instanceName)
	d.Set("name", backup.Note)
	d.Set("state", backup.Status)
	d.Set("size", backup.Size)
	d.Set("type", backup.Type)
	d.Set("created_at", backup.CreatedAt)

	return nil
}

func snapshotInfo(d *schema.ResourceData, meta interface{}) (*ah.Backup, string, error) {
	client := meta.(*ah.APIClient)
	options := &ah.ListOptions{
		Filters: []ah.FilterInterface{
			&ah.EqFilter{
				Keys:  []string{"id"},
				Value: d.Id(),
			},
		},
	}
	instanceBackups, err := client.Backups.List(context.Background(), options)
	if err != nil {
		return nil, "", err
	}

	if len(instanceBackups) != 1 {
		return nil, "", ah.ErrResourceNotFound
	}

	instanceBackup := instanceBackups[0]
	return &instanceBackup.Backups[0], instanceBackup.InstanceName, nil

}

func resourceAHCloudServerSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	if d.HasChange("name") {

		updateRequest := &ah.BackUpUpdateRequest{
			Note: d.Get("name").(string),
		}

		if _, err := client.Backups.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"Error changing backup name (%s): %s", d.Id(), err)
		}
	}

	return resourceAHCloudServerSnapshotRead(d, meta)
}

func resourceAHCloudServerSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if _, err := client.Backups.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error deleting backup (%s): %s", d.Id(), err)
	}
	if err := waitForBackupDestroy(d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for backup (%s) to become destroyed: %s", d.Id(), err)
	}
	return nil
}

func waitForBackupReady(instanceID, actionID string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		action, err := client.Instances.ActionInfo(context.Background(), instanceID, actionID)
		if err != nil {
			log.Printf("Error on waitForBackupReady: %v", err)
			return nil, "", err
		}
		return action.ID, action.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:                     2 * time.Second,
		Pending:                   []string{"queued", "running"},
		Refresh:                   stateRefreshFunc,
		Target:                    []string{"success"},
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for backup active status: %v", err)
	}
	return nil

}

func waitForBackupDestroy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		backup, err := client.Backups.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil {
			log.Printf("Error on waitForBackupDestroy: %v", err)
			return nil, "", err
		}

		return backup.ID, backup.Status, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"pending_delete"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for backup to be destroyed: %s", err)
	}

	return nil
}
