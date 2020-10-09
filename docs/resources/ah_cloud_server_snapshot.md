# AH Cloud Server Snapshot Resource

Provides an AdvancedHosting Cloud Server Snapshot resource which can be used to create a Cloud Server Snapshot.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "ubuntu-20-04-x64"
  name = "Sample server"
  datacenter = "ams1"
  product = "start-m"
}

resource "ah_cloud_server_snapshot" "example-snapshot" {
  cloud_server_id = ah_cloud_server.example.id
  name = "example-snapshot-1"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_server_id` - (Required) Cloud Server ID to create a Snapshot from.
* `name` - (Optional) Name of the snapshot. If not set, the snapshot name is assigned automatically based on date and time of snapshot creation.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of the Snapshot.
* `cloud_server_name` - Cloud Server name the Snapshot was created for.
* `state` - Current state of the Snapshot.
* `size` - Snapshot size, in GB
* `type` - Type. Can be `snapshot` (for manual snapshots) or `backup` (for automatic backups)
* `created_at` - Creation datetime of the Snapshot.
