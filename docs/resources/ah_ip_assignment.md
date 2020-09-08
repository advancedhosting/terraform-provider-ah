---
layout: "ah"
page_title: "Advahced Hosting: ah_ip_assignment"
sidebar_current: "docs-resource-ah-ip-assignment"
description: |-
  Manages Advanced Hosting Cloud Server IP address assignments.
---

# ah_ip_assignment

Provides an Advanced Hosting IP Assignment resource to assign an IP address to a Cloud Server. Can be done either using this resource or in `ips` argument of `ah_cloud_server`.

If IP address is not set using `ah_ip_assignment` and not provided in the list of IP addresses defined in the `ips` argument of `ah_cloud_server`, IP address will be unassigned from the Server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "ubuntu-20-04-x64"
  name = "Sample server"
  datacenter = "ams1"
  product = "start-m"
}

resource "ah_ip" "example" {
  type = "public"
  datacenter = "ams1"
}

resource "ah_ip_assignment" "example" {
  cloud_server_id = ah_cloud_server.example.id
  ip_address = ah_ip.example.id
}

```

## Argument Reference

The following arguments are supported:

* `cloud_server_id` - (Required) Cloud Server ID to assign IP addresses to.
* `ip_address` - (Required) IP address ID or IP address value to assign.
* `primary` - (Optional) - Boolean for the Primary IP flag. Only one of assignments can have this flag set to true. Default value is false.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Unique ID of the IP Address Assignment.