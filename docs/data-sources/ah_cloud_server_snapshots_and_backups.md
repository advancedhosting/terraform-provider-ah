# AH Cloud Server Snapshots And Backups Data Source

Get information about AdvancedHosting Cloud Server Snapshots and Automatic Backups.

## Example Usage

Get the Cloud Server Snapshot or Automatic Backup by ID:

```hcl
data "ah_cloud_server_snapshot_and_backups" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get the Snapshot or Automatic Backup list for specific Cloud Server:

```hcl
data "ah_cloud_server_snapshot_and_backups" "example" {
  filter {
    key = "cloud_server_id"
    values = ["123"]
  }
}
```

Get a list of Snapshots created for specific Cloud Server, sorted by creation date:

```hcl
data "ah_cloud_server_snapshot_and_backups" "example" {
  filter {
    key = "cloud_server_id"
    values = ["123"]
  }
  filter {
    key = "type"
    values = ["snapshot"]
  }
  sort {
    key = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter`: (Optional) Filter the results by specified key and value. The structure of the block is documented below.
* `sort` - (Optional) Sort the results by specified key and direction. The structure of the block is documented below.

---

The `filter` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `cloud_server_id`,  `cloud_server_name`, `state`, `size`, `type`.
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `cloud_server_id`,  `cloud_server_name`, `state`, `size`, `type`,`created_at`.
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `snapshots_and_backups` - A list of snapshots and backups that satisfy the search criteria.
    * `id` - ID of the Snapshot.
    * `name` - Name of the snapshot. 
    * `cloud_server_id` - Cloud Server ID a Snapshot was created from.
    * `cloud_server_name` - Cloud Server Name snapshot was created for.
    * `cloud_server_deleted` - Boolen flag indicating whether the original Cloud Server was deleted.
    * `state` - Current state of the Snapshot.
    * `size` - Snapshot size, in GB
    * `type` - Type. Can be `snapshot` (for manual snapshots) or `backup` (for automatic backups).
    * `created_at` - Creation datetime of the Snapshot.
