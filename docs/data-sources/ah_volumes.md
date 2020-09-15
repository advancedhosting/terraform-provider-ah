---
layout: "ah"
page_title: "Advahced Hosting: ah_volumes"
sidebar_current: "docs-data-source-ah-volumes"
description: |-
  Get information about Advanced Hosting Volumes.
---

# ah_volumes

Get information about Advanced Hosting Volumes.

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
    values = ["hdd-level2-ams1", "hdd-level3-ams1"]
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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product_id`, `size`, `file_system`.
<!-- * `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product`, `size`, `file_system`, `cloud_server_id`, `created_at` TODO add cloud_server_id filter WCS-3584-->
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product_id`, `size`, `file_system`, `created_at`
<!-- * `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `product`, `size`, `file_system`, `cloud_server_id`, `created_at` TODO add cloud_server_id sorting WCS-3584--> 
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