---
layout: "ah"
page_title: "Advahced Hosting: ah_cloud_servers"
sidebar_current: "docs-data-source-ah-cloud-servers"
description: |-
  Get information about Advanced Hosting Cloud Servers.
---

# ah_cloud_servers

Get information about Advanced Hosting Cloud Servers.

## Example Usage

Get the Cloud Server by ID:

```hcl
data "ah_cloud_servers" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of active Cloud Servers from AMS1 and ASH1 datacenter, sorted by creation date, desc:

```hcl
data "ah_cloud_servers" "example" {
  filter {
    key = "state"
    values = ["active"]
  }
  filter {
    key = "datacenter"
    values = ["ams1", "ash1"]
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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `state`, `current_action`,  `vcpu`, `ram`, `disk`, `created_at`, `ip_address`, `private_network_id`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `state`, `current_action`,  `vcpu`, `ram`, `disk`, `created_at`, `ip_address`, `private_network_id`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `cloud_servers` - A list of Cloud Servers that satisfy the search criteria.
  * `id` -  ID of the server.
  * `name` - Name of the Cloud Server.
  * `datacenter` - Datacenter Slug of the Cloud Server.
  * `product` - Product Slug that indentifies the product type of Cloud Server.
  * `state` - Current state of the Cloud Server.
  * `current_action` - Current action is being performed on the server.
  * `vcpu` - Number of vCPUs on the server.
  * `ram` - RAM of the server in GiB.
  * `disk` - Disk size of the server in GB.
  * `created_at` - Creation timestamp of the Server.
  * `image` - Cloud Server Image ID or Snapshot / Auto Backup ID the server was created from.
  * `backups` - Boolean indicating whether backups are enabled for the server.
  * `use_password` - Boolean indicating whether server was created with a password generated.
  * `ssh_keys` - Array of SSH fingerprints the server was created with.
  * `ips` - Array of Public and Anycast IP addresses assigned to the server. The structure of the block is documented below.
  * `volumes` - Array of Volume IDs attached to the server.
  * `private_networks` - Array of Private Networks connected to the server. The structure of the block is documented below.
  * `firewall_rules` - Array of Firewall Rules applied to the server. The structure of the block is documented below.

---

The `ips` block contains:
* `ip_address` - Public IP address assigned to the Cloud Server.
* `type` - IP address type. Can be either `public` or `anycast`.
* `primary` - Boolean indicating a Primary IP flag.
* `reverse_dns` - Reverse DNS assigned to the IP address.
* `assignment_id` - ID of the IP Address Assignment.

The `private_networks` block contains:

* `id` - Private Network ID.
* `ip` - Private network IP address of the Cloud Server within the network.

The `firewall_rules` block supports:

* `type` - Type of the rule. Can be either `inbound` or `outbound`.
* `action`- Type of action for the rule. Can be either `accept` or `drop`.
* `traffic_type` - Type of the traffic to apply the rule to. Can be one of: `all`, `icmp`, `tcp`, `udp`.
* `ip_range` - IP address range for the rule in CIDR format (e.g. '10.4.4.4/24').
* `ports` - List of Ports for the rule. Can be a port number (80) or a port range (1-65535).

