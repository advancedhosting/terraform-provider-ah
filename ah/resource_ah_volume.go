package ah

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAHVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHVolumeCreate,
		Read:   resourceAHVolumeRead,
		Update: resourceAHVolumeUpdate,
		Delete: resourceAHVolumeDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "New volume",
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_system": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ext4",
				ValidateFunc: validation.StringInSlice([]string{"ext4", "btrfs", "xfs"}, false),
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("size", func(old, new, meta interface{}) error {
				if new.(int) < old.(int) {
					return fmt.Errorf("New size value must be greater than old value %d", old.(int))
				}
				return nil
			}),
		),
	}
}

func resourceAHVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	request := &ah.VolumeCreateRequest{
		Name:       d.Get("name").(string),
		Size:       d.Get("size").(int),
		ProductID:  d.Get("product").(string),
		FileSystem: d.Get("file_system").(string),
	}

	volume, err := client.Volumes.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}

	d.SetId(volume.ID)

	if err := waitForVolumeState([]string{"creating"}, []string{"ready"}, d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become ready: %s", d.Id(), err)
	}

	return resourceAHVolumeRead(d, meta)

}

func resourceAHVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	volume, err := client.Volumes.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", volume.Name)
	d.Set("product", volume.ProductID)
	d.Set("size", volume.Size)
	d.Set("file_system", volume.FileSystem)

	return nil
}

func resourceAHVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	d.Partial(true)

	if d.HasChange("name") {

		updateRequest := &ah.VolumeUpdateRequest{
			Name: d.Get("name").(string),
		}

		if _, err := client.Volumes.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"Error changing volume name (%s): %s", d.Id(), err)
		}
		d.SetPartial("name")
	}

	if d.HasChange("size") {
		if _, err := client.Volumes.Resize(context.Background(), d.Id(), d.Get("size").(int)); err != nil {
			return fmt.Errorf(
				"Error resizing volume (%s): %s", d.Id(), err)
		}
		if err := waitForVolumeState([]string{"resizing"}, []string{"ready"}, d, meta); err != nil {
			return fmt.Errorf(
				"Error waiting for volume (%s) to become ready: %s", d.Id(), err)
		}
		d.SetPartial("size")
	}

	return resourceAHVolumeRead(d, meta)
}

func resourceAHVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.Volumes.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error deleting volume (%s): %s", d.Id(), err)
	}
	if err := waitForVolumeDestroy(d, meta); err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become destroyed: %s", d.Id(), err)
	}
	return nil
}

func waitForVolumeState(pending, target []string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		volume, err := client.Volumes.Get(context.Background(), d.Id())
		if err != nil || volume == nil {
			log.Printf("Error on waitForVolumeState: %v", err)
			return nil, "", err
		}
		return volume.ID, volume.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      2 * time.Second,
		Pending:    pending,
		Refresh:    stateRefreshFunc,
		Target:     target,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume status %v: %v", target, err)
	}

	return nil

}

func waitForVolumeDestroy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		instance, err := client.Volumes.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil || instance == nil {
			log.Printf("Error on waitForVolumeDestroy: %v", err)
			return nil, "", err
		}

		return instance.ID, instance.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"deleting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume to be destroyed: %s", err)
	}

	return nil
}
