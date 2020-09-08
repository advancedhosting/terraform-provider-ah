---
layout: "ah"
page_title: "Advahced Hosting: ah_ips"
sidebar_current: "docs-data-source-ah-ips"
description: |-
  Get information about Advanced Hosting Public and Anycast IP addresses.
---

# ah_ips

Get information about Advanced Hosting Public and Anycast IP addresses.

## Example Usage

Get the IP address by ID:

```hcl
data "ah_ips" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get the IP address by IP address value:

```hcl
data "ah_ips" "example" {
  filter {
    key = "ip"
    values = ["1.2.3.4"]
  }
}
```

Get a list of public IP addresses assigned to a Cloud Server, sorted by the creation date:

```hcl
data "ah_ips" "example" {
  filter {
    key = "cloud_server_id"
    values = ["123"]
  }
  filter {
      key = "type"
      values = ["public"]
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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `reverse_dns`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `ip_address`, `reverse_dns`, `created_at`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `ips` - A list of IP addresses that satisfy the search criteria.
    * `id` - ID of IP address.
    * `ip_address` - IP address value.
    * `type` - IP address type. Can be either `public` or `anycast`.
    * `datacenter` - Datacenter Slug where IP addresses is allocated (returned only if `type="public"`). 
    * `reverse_dns` - Reverse DNS assigned to the IP address.
    * `cloud_server_ids` - List of Cloud Server IDs the IP addresses is assigned to.
    * `created_at` - Creation datetime of the IP address.
    * `primary` - Boolean for the Primary IP flag. Only IPs of `public` type will have this flag, can contain a value only if IP is assigned to a server.