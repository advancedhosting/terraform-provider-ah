package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAHVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHVolumeAttachmenCreate,
		Read:   resourceAHVolumeAttachmenRead,
		Delete: resourceAHVolumeAttachmenDelete,
		Schema: map[string]*schema.Schema{
			"cloud_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAHVolumeAttachmenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	instanceID := d.Get("cloud_server_id").(string)

	volumeID := d.Get("volume_id").(string)

	action, err := client.Instances.AttachVolume(context.Background(), instanceID, volumeID)

	if err != nil {
		return err
	}

	if err := waitForVolumeState(volumeID, []string{"attaching"}, []string{"attached"}, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become attached: %s", volumeID, err)
	}

	d.SetId(action.ID)

	return resourceAHVolumeAttachmenRead(d, meta)

}

func resourceAHVolumeAttachmenRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ah.APIClient)
	volumeID := d.Get("volume_id").(string)

	volume, err := client.Volumes.Get(context.Background(), volumeID)
	if err != nil {
		return err
	}
	if volume.Instance == nil {
		return fmt.Errorf("Error reading for volume attaching info: volume %s is not attached", volumeID)
	}
	d.Set("cloud_server_id", volume.Instance.ID)
	d.Set("volume_id", volume.ID)
	d.Set("state", volume.State)
	return nil
}

func resourceAHVolumeAttachmenDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ah.APIClient)

	instanceID := d.Get("cloud_server_id").(string)

	volumeID := d.Get("volume_id").(string)

	if _, err := client.Instances.DetachVolume(context.Background(), instanceID, volumeID); err != nil {
		return err
	}

	if err := waitForVolumeState(volumeID, []string{"detaching"}, []string{"ready"}, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become detached: %s", volumeID, err)
	}

	return nil
}
