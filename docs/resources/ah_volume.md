# AH Volume Resource

Provides an Advanced Hosting Volume resource.

## Example Usage

```hcl

resource "ah_volume" "example" {
  name = "Volume Name"
  product = "hdd-l2-ash1"
  file_system = "ext4"
  size = "20"
}

```

## Argument Reference

The following arguments are supported:
* `name` - (Optional) Volume name
* `product` - (Required) Volume Product ID or Slug that indentifies the desired product type of Volume. Changing this erases and recreates the volume. See the [list of available volume products](https://websa.advancedhosting.com/slugs).
* `size` - (Optional) Desired volume size in GB. Changing allowed to a greater value only. Changing this increases the volume size, data is preserved. Required unless `origin_volume_id` is set.
* `file_system` - (Optional) File system formatting option. Can be one of: `ext4`, `btrfs`, `xfs`, or empty. If empty, volume is not formatted. Default value is `ext4`. Changing this erases and recreates the volume.
* `origin_volume_id` - (Optional) ID of the volume to copy from.  Changing this erases and recreates the volume.  If this argument is set, `size` is ignored.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of the Volume
* `state` - Current state of the Volume.
* `created_at` - Creation datetime of the Volume.