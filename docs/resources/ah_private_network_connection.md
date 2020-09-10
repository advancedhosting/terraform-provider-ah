---
layout: "ah"
page_title: "Advahced Hosting: ah_private_network_connection"
sidebar_current: "docs-resource-ah-private-network-connection"
description: |-
  Manages Advanced Hosting Cloud Server Private Network connection.
---

# ah_private_network_connection

Provides an Advanced Hosting Private Network Connection resource to connect a Cloud Server to a Private Netowrk.


## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "f0438a4b-7c4a-4a63-a593-8e619ec63d16"
  name = "Sample server"
  datacenter = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
  product = "df42a96b-b381-412c-a605-d66d7bf081af"
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