package ah

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAHCloudServerSnapshotsAndBackups_Basic(t *testing.T) {

	resourcesConfig := `
	resource "ah_cloud_server" "web" {
	  count = 3
	  name = "test_${count.index}"
	  datacenter = "ams1"
	  image = "centos7-64"
	  product = "start-xs"
	}

	resource "ah_cloud_server_snapshot" "test" {
	  cloud_server_id = ah_cloud_server.web.0.id
	  name = "Test Name"
	}
	
	resource "ah_cloud_server_snapshot" "test2" {
		cloud_server_id = ah_cloud_server.web.1.id
		name = "New Name"
	}`

	datasourceConfig := `
	data "ah_cloud_server_snapshot_and_backups" "test" {
	  filter {
		key = "cloud_server_id"
		values = [ah_cloud_server.web.0.id]
	  }
	  sort {
		key = "created_at"
		direction = "desc"
	  }
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.#", "1"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.id", "ah_cloud_server_snapshot.test", "id"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.name", "ah_cloud_server_snapshot.test", "name"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.cloud_server_id", "ah_cloud_server_snapshot.test", "cloud_server_id"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.cloud_server_name", "ah_cloud_server_snapshot.test", "cloud_server_name"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.state", "ah_cloud_server_snapshot.test", "state"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.size", "ah_cloud_server_snapshot.test", "size"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.type", "ah_cloud_server_snapshot.test", "type"),
					resource.TestCheckResourceAttrPair("data.ah_cloud_server_snapshot_and_backups.test", "snapshots_and_backups.0.created_at", "ah_cloud_server_snapshot.test", "created_at"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
