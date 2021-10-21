package ah

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use plan instead",
			},
			"plan": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"product"},
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"file_system": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ext4",
				ValidateFunc: validation.StringInSlice([]string{"ext4", "btrfs", "xfs", ""}, false),
			},
			"origin_volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("size", func(ctx context.Context, old, new, meta interface{}) error {
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
	name := d.Get("name").(string)

	var planAttr string

	plan, planOk := d.GetOk("plan")
	product, productOk := d.GetOk("product")
	if !planOk && !productOk {
		return errors.New("one of plan or product must be configured")
	}

	if planOk {
		planAttr = plan.(string)
	} else {
		planAttr = product.(string)
	}

	if attr, ok := d.GetOk("origin_volume_id"); ok {
		request := &ah.VolumeCopyActionRequest{
			Name: name,
		}

		if planID, err := strconv.Atoi(planAttr); err != nil {
			request.PlanSlug = planAttr
		} else {
			request.PlanID = planID
		}

		originVolumeID := attr.(string)
		action, err := client.Volumes.Copy(context.Background(), originVolumeID, request)
		if err != nil {
			return fmt.Errorf("error creating volume from origin: %s", err)
		}
		if err := waitForActionCopyReady(originVolumeID, action.ID, d, meta); err != nil {
			return err
		}
		action, err = client.Volumes.ActionInfo(context.Background(), originVolumeID, action.ID)
		if err != nil {
			return err
		}
		d.SetId(action.ResultParams.CopiedVolumeID)
	} else {

		request := &ah.VolumeCreateRequest{
			Name:       name,
			Size:       d.Get("size").(int),
			FileSystem: d.Get("file_system").(string),
		}

		if planID, err := strconv.Atoi(planAttr); err != nil {
			request.PlanSlug = planAttr
		} else {
			request.PlanID = planID
		}

		volume, err := client.Volumes.Create(context.Background(), request)

		if err != nil {
			return fmt.Errorf("Error creating volume: %s", err)
		}
		d.SetId(volume.ID)
	}

	if err := waitForVolumeState(d.Id(), []string{"creating"}, []string{"ready"}, d, meta); err != nil {
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
	d.Set("size", volume.Size)
	d.Set("file_system", volume.FileSystem)
	d.Set("state", volume.State)
	d.Set("created_at", volume.CreatedAt)

	return nil
}

func resourceAHVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	if d.HasChange("name") {

		updateRequest := &ah.VolumeUpdateRequest{
			Name: d.Get("name").(string),
		}

		if _, err := client.Volumes.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf(
				"error changing volume name (%s): %s", d.Id(), err)
		}

	}

	if d.HasChange("size") {
		if _, err := client.Volumes.Resize(context.Background(), d.Id(), d.Get("size").(int)); err != nil {
			return fmt.Errorf(
				"Error resizing volume (%s): %s", d.Id(), err)
		}
		if err := waitForVolumeState(d.Id(), []string{"resizing"}, []string{"ready", "attached"}, d, meta); err != nil {
			return fmt.Errorf(
				"Error waiting for volume (%s) to become ready: %s", d.Id(), err)
		}

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

func waitForVolumeState(volumeID string, pending, target []string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		volume, err := client.Volumes.Get(context.Background(), volumeID)
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
			"error waiting for volume status %v: %v", target, err)
	}

	return nil

}

func waitForVolumeDestroy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		volume, err := client.Volumes.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil || volume == nil {
			log.Printf("Error on waitForVolumeDestroy: %v", err)
			return nil, "", err
		}

		return volume.ID, volume.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"deleting", "detaching"},
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

func waitForActionCopyReady(VolumeID, actionID string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		action, err := client.Volumes.ActionInfo(context.Background(), VolumeID, actionID)
		if err != nil {
			log.Printf("Error on waitForActionCopyReady: %v", err)
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
		ContinuousTargetOccurence: 2,
	}
	_, err := stateChangeConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume coping success status: %v", err)
	}
	return nil

}
