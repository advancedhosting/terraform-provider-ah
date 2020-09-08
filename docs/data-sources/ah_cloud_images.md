---
layout: "ah"
page_title: "Advahced Hosting: ah_cloud_images"
sidebar_current: "docs-data-source-ah-cloud-images"
description: |-
  Get information about Advanced Hosting Cloud Images available for server creation.
---

# ah_cloud_images

Get information about Advanced Hosting Cloud Images available for server creation.

## Example Usage

Get the Image by ID:

```hcl
data "ah_cloud_images" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of Ubuntu and Debian x64 images, sorted by version, desc:

```hcl
data "ah_cloud_images" "example" {
  filter {
    key = "distribution"
    values = ["Ubuntu", "Debian"]
  }
  filter {
    key = "architecture"
    values = ["x64"]
  }
  sort {
    key = "version"
    direction = "desc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter`: (Optional) Filter the results by specified key and value. The structure of the block is documented below.
* `sort` - (Optional) Sort the results by specified key and direction. The structure of the block is documented below.

---
* TODO либо убрать линию, либо добавить после filter
* TODO сейчас таких фильтраций нет, нужно ли их добавлять на стороне апи или на стороне TF ?
* TODO нужно добавить data source и для одного ресурса

The `filter` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `distribution`,  `version`, `architecture`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `distribution`,  `version`, `architecture`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `images` - A list of Images that satisfy the search criteria.
  * `id` - ID of the Image.
  * `name` - Name of the Image.
  * `distribution` - Name of the Image Distribution.
  * `version` - Distribution version.
  * `architecture` - Distribution architecture.
