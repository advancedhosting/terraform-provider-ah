package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/crypto/ssh"
)

func resourceAHSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAHSSHKeyCreate,
		Read:   resourceAHSSHKeyRead,
		Update: resourceAHSSHKeyUpdate,
		Delete: resourceAHSSHKeyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func resourceAHSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	publicKey := d.Get("public_key").(string)
	_, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return fmt.Errorf("Invalid public_key: %v", err)
	}

	request := &ah.SSHKeyCreateRequest{
		PublicKey: publicKey,
	}

	if attr, ok := d.GetOk("name"); ok {
		request.Name = attr.(string)
	} else {
		request.Name = comment
	}

	sshKey, err := client.SSHKeys.Create(context.Background(), request)

	if err != nil {
		return fmt.Errorf("Error creating ssh key: %s", err)
	}

	d.SetId(sshKey.ID)

	return resourceAHSSHKeyRead(d, meta)

}

func resourceAHSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	sshKey, err := client.SSHKeys.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", sshKey.Name)
	d.Set("public_key", sshKey.PublicKey)
	d.Set("fingerprint", sshKey.Fingerprint)
	d.Set("created_at", sshKey.CreatedAt)

	return nil
}

func resourceAHSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	updateRequest := &ah.SSHKeyUpdateRequest{}

	if d.HasChange("name") {
		updateRequest.Name = d.Get("name").(string)
	}

	if d.HasChange("public_key") {
		updateRequest.PublicKey = d.Get("public_key").(string)
	}

	if updateRequest.Name == "" || updateRequest.PublicKey == "" {
		if _, err := client.SSHKeys.Update(context.Background(), d.Id(), updateRequest); err != nil {
			return fmt.Errorf("Error updating ssh_key (%s): %s", d.Id(), err)
		}
	}

	return resourceAHSSHKeyRead(d, meta)
}

func resourceAHSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	if err := client.SSHKeys.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf(
			"Error deleting ssh key (%s): %s", d.Id(), err)
	}
	return nil
}
