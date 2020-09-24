---
layout: "ah"
page_title: "Advahced Hosting: ah_cloud_server"
sidebar_current: "docs-resource-ah-cloud-server"
description: |-
    Manages Advanced Hosting Cloud Servers.
---

# ah_cloud_server

Provides an Advanced Hosting Cloud Server resource. This can be used to create, modify, delete Cloud Servers, manage IP address assignments, volume attachments, private network connections and firewall rules for the server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
  name = "Sample server"
  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
  product = "df42a96b-b381-412c-a605-d66d7bf081af"
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The Cloud Server image ID or Slug of the desired image for the server OR a Cloud Server Snapshot / Auto Backup ID. Changing this creates a new server.
* `name` - (Required) Name for the Cloud Server.
* `datacenter` - (Required) Datacenter ID or Slug to start the Cloud Server in.
* `product` - (Required) Cloud Server Product ID or Slug that indentifies the desired product type of the Cloud Server. Changing this resizes the existing server.
* `backups` - (Optional) Boolean to enable or disable backups. Defaults to false.
* `use_password` - (Optional) Boolean defining if password should be generated for the server and sent by email. Defaults to true.
* `ssh_keys` - (Optional) Array of SSH keys IDs to enable in
   the format `[12345, 8595645]`.
* `create_public_ip_address` - (Optional) Boolean defining if a new public IP address should be created for the server. This public IP address will become a primary IP address for the Cloud Server. Defaults to true.
---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of the server.
* `state` - Current state of the Cloud Server.
* `vcpu` - Number of vCPUs on the the Cloud Server.
* `ram` - RAM of the Cloud Server in MiB.
* `disk` - Disk size of the Cloud Server in GB.
* `created_at` - Creation datetime of the Cloud Server.
* `ips` -  Array of IP address blocks to be assigned to the Cloud Server.

---

The `ips` block contains:
* `ip_address` - Public IP address assigned to the Cloud Server.
* `type` - IP address type. Can be either `public` or `anycast`.
* `primary` - Boolean indicating a Primary IP flag.
* `reverse_dns` - Reverse DNS assigned to the IP address.
* `assignment_id` - ID of the IP Address Assignment.
