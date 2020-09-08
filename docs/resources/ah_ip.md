---
layout: "ah"
page_title: "Advahced Hosting: ah_ip"
sidebar_current: "docs-resource-ah-ip"
description: |-
    Manages Advanced Hosting Cloud Server IP addresses.
---

# ah_ip

Provides an Advanced Hosting IP address resource to represent a publicly-accessible static public and anycast IP addresses that can be mapped to your servers.


## Example Usage

```hcl

resource "ah_ip" "example" {
  type = "public"
  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
}

```

## Argument Reference

The following arguments are supported:

* `type` - (Required) IP type. Can be either `public` or `anycast`.
* `datacenter` - (Optional) Datacenter ID to create IP addresses. Required if `type="public"`, ignored if `type="anycast"`.
* `reverse_dns` - (Optional) Reverse DNS to be assigned to the IP address. If not specified, will be automatically generated.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of IP address.
* `ip_address` - IP address value.
* `created_at` - Creation datetime of the IP address.