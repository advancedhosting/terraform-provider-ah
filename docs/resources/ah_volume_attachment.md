# AH Volume Attachment Resource

Provides an Advanced Hosting Volume Attachment resource to connect a Volume to the Cloud Server. Can be done either using this resource or in `volumes` argument of `ah_cloud_server`.

If Volume is not set using `ah_volume_attachment` and not provided in the list of Volumes defined in the `volumes` argument of `ah_cloud_server`, Volume will be detached from the Cloud Server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "ubuntu-20-04-x64"
  name = "Sample server"
  datacenter = "ams1"
  product = "start-m"
}

resource "ah_volume" "example" {
  name = "Volume Name"
  product = "hdd-level2-ams1"
  file_system = "ext4"
  size = "20"
}

resource "ah_volume_attachment" "example" {
  cloud_server_id = ah_cloud_server.example.id
  volume_id = ah_volume.example.id
}

```

## Argument Reference

The following arguments are supported:

* `cloud_server_id` - (Required) Cloud Server ID to attach a Volume to.
* `volume_id` - (Required) Volume ID to attach to a Cloud Server.
---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Unique ID of the Volume Attachment.
* `state` - Current state of attachment.