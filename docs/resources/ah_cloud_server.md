---
layout: "ah"
page_title: "Advahced Hosting: ah_cloud_server"
sidebar_current: "docs-resource-ah-cloud-server"
description: |-
    Manages Advanced Hosting Cloud Servers.
---

# ah_cloud_server

Provides an Advanced Hosting Cloud Server resource. This can be used to create, modify, delete Cloud Servers, manage IP address assignments, volume attachments, private network connections, and firewall rules for the server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  image = "ubuntu-20-04-x64"
  name = "Sample server"
  datacenter = "ams1"
  product = "start-m"
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) Cloud Server Image ID or Slug of the desired image for the server OR a Cloud Server Snapshot / Auto Backup ID. Changing this creates a new server.
* `name` - (Required) Name for the Cloud Server.
* `datacenter` - (Required) Datacenter ID or Slug to start Cloud Server in.
* `product` - (Required) Cloud Server Product ID or Slug that indentifies the desired product type of Cloud Server. Changing this resizes the existing server.
* `backups` - (Optional) Boolean to enable or disable backups. Defaults to false.
* `use_password` - (Optional) Boolean defining if password should be generated for the server and sent by email. Defaults to true.
* `ssh_keys` - (Optional) Array of SSH IDs or fingerprints to enable in
   the format `[12345, 7e:ac:a8:45:83:e3:58:f5:3a:9f:dd:16:63:dc:fb:1e]`. Fingerprints can be found in the 'SSH keys' section of the panel.
* `create_public_ip_address` - (Optional) Boolean defining if a new public IP address should be created for the server. This public IP address will become a primary IP address for the server. Defaults to true.
* `ips` - (Optional) Array of IP address blocks to be assigned to the server. The structure of the block is documented below.
* TODO сейчас мы не моэем приконектить при создании через api. Делаем это руками через TF ? Да, делаем на стороне TF
* `volumes` - (Optional) Array of Volume IDs to be attached to the server in the format `[12345, 123456]`.
* `private_networks` - (Optional) Array of Private Network blocks to be connected to the server. The structure of the block is documented below.
* `firewall_rules` - (Optional) Array of Firewall Rule blocks to be applied to the server. The structure of the block is documented below.

---

The `ips` block supports:

* `ip_address` - (Required) Public or Anycast IP address (ID or IP address value) to be assigned to the Cloud Server.
* `primary` - (Optional) - Boolean indicating the Primary IP flag. Can be set for Public IP addresses only. Only one of the values can have this flag set to true. Default value is false.

The `private_networks` block supports:

* `id` - (Required) Private Network ID.
* `ip_address` - (Optional) Private Network IP address of the Cloud Server within the network. If not set, IP is assigned automatically.

The `firewall_rules` block supports:

* `type` - (Required) Type of the rule. Can be either `inbound` or `outbound`.
* `action`- (Required) Type of action for the rule. Can be either `accept` or `drop`.
* `traffic_type` - (Optional) Type of the traffic to apply the rule to. Can be one of: `all`, `icmp`, `tcp`, `udp`.
* `ip_range` - (Optional) IP address range for the rule in CIDR format (e.g. '10.4.4.4/24'). Default value is empty which sets the rule to apply all IP addresses.
* `ports` - (Optional) List of Ports for the rule. Can be a port number (80) or a port range (1-65535). Default value is empty which sets the rule to apply all ports.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of the server.
* `state` - Current state of the Cloud Server.
* `current_action` - Current action that is being performed on the server.
* `vcpu` - Number of vCPUs on the server.
* `ram` - RAM of the server in GiB.
* `disk` - Disk size of the server in GB.
* `created_at` - Creation datetime of the Server.

---

The `ips` block additionally contains:

* `type` - IP address type. Can be either `public` or `anycast`.
* TODO при создании мы нигде это не указываем
* `reverse_dns` - Reverse DNS assigned to the IP address.
* `assignment_id` - ID of the IP Address Assignment.
