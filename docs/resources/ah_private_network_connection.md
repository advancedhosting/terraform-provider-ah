---
layout: "ah"
page_title: "Advahced Hosting: ah_private_network_connection"
sidebar_current: "docs-resource-ah-private-network-connection"
description: |-
  Manages Advanced Hosting Cloud Server Private Network connection.
---

# ah_private_network_connection

Provides an Advanced Hosting Private Network Connection resource to connect a Cloud Server to a Private Netowrk. Can be done either using this resource or in `private_networks` argument of `ah_cloud_server`.

If Private Network is not set using `ah_private_network_connection` and not provided in the list of Private Netwroks defined in the `private_networks` argument of `ah_cloud_server`, Private Netwrok will be disconnected from the Server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "ubuntu-20-04-x64"
  name = "Sample server"
  datacenter = "ams1"
  product = "start-m"
}

resource "ah_private_network" "example" {
  ip_range = "10.0.0.0/24"
  name = "New Private Network"
}

resource "ah_private_network_connection" "example" {
  cloud_server_id = ah_cloud_server.example.id
  private_network_id = ah_private_network.example.id
}

```

## Argument Reference

The following arguments are supported:

* `cloud_server_id` - (Required) Cloud Server ID to connect to a Private Network.
* `private_network_id` - (Required) Private Network ID to connect a server to.
* `ip_address` - (Optional) Private network IP address of the Cloud Server within the network. If not set, IP is assigned automatically.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Unique ID of the Private Network Connection.
