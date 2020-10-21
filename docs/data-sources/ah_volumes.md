# AH Volumes

Get information about AdvancedHosting Volumes.

## Example Usage

Get the Volume by ID:

```hcl
data "ah_volumes" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of HDD Volumes attached to a Cloud Server sorted by creation date:

```hcl
data "ah_volumes" "example" {
  filter {
    key = "cloud_server_id"
    values = ["123"]
  }
  filter {
    key = "product"
    values = ["hdd-l2-ash1", "hdd3-ash1"]
  }
  sort {
    key = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter`: (Optional) Filter the results by specified key and value. The structure of the block is documented below.
* `sort` - (Optional) Sort the results by specified key and direction. The structure of the block is documented below.

---

The `filter` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product_id`, `size`, `file_system`, `cloud_server_id`.
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product_id`, `size`, `file_system`, `created_at`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `volumes` - A list of Volumes that satisfy the search criteria.
    * `id` - ID of the Volume.
    * `name` - Volume name.
    * `state` - Current state of the Volume.
    * `product` - Product Slug that indentifies the product type of Volume. 
    * `size` - Volume size in GB.
    * `file_system` - File system formatting option selected on volume creation. Can be one of: `ext4`, `btrfs`, `xfs`, or empty.
    * `cloud_server_id` - Cloud Server ID the Volumes is attached to, if attached.
    * `created_at` - Creation datetime of the Volume.