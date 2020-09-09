---
layout: "ah"
page_title: "Advahced Hosting: ah_private_networks"
sidebar_current: "docs-data-source-ah-private-networks"
description: |-
  Get information about Advanced Hosting Private Networks.
---

# ah_private_networks

Get information about Advanced Hosting Private Networks.

## Example Usage

Get a Private Network by ID:

```hcl
data "ah_private_networks" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```


Get a list of active Private Networks connected to a set of Cloud Servers, sorted by the creation date:

```hcl
data "ah_private_networks" "example" {
    filter {
        key = "cloud_server_id"
        values = ["123", "456"]
    }
    filter {
        key = "state"
        values = ["active"]
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
<!-- * `key` - (Required) Filter the results by specified key. Can be one of: `id`, `ip_range`, `name`,  `state`, `cloud_server_id`, `created_at` -->
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `ip_range`, `name`,  `cloud_server_id`, `state`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
<!-- * `key` - (Required) Filter the results by specified key. Can be one of: `id`, `ip_range`, `name`,  `state`, `cloud_server_id`, `created_at` -->
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `ip_range`, `name`,  `state`, `created_at`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `private_networks` - A list of Private Networks that satisfy the search criteria.
    * `id` -  ID of the Private Network.
    * `ip_range` - Private Network IP range in CIDR format.
    * `name` - (Optional) Name of the Private Network.
    * `state` - Current state of the Private Network.
    * `cloud_servers` - List of Cloud Servers the Private Network is connected to. The structure of the block is documented below.
    * `created_at` - Creation datetime of the Private Network.

---

The `cloud_servers` block contains:

* `id` - Cloud Server ID.
* `ip` - Private network IP address of the Cloud Server within the network.