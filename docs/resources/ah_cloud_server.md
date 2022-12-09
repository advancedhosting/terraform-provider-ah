# AH Cloud Server Resource

Provides an Advanced Hosting Cloud Server resource. This can be used to create, modify, delete Cloud Servers, manage IP address assignments, volume attachments, private network connections and firewall rules for the Cloud Server.

## Example Usage

```hcl
resource "ah_cloud_server" "example" {
  name = "Sample server"
  datacenter = "ams1"
  image = "centos-7-x64"
  plan = "start-xs"
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The Cloud Server image ID or Slug of the desired image for the server OR a Cloud Server Snapshot / Auto Backup ID. Changing this creates a new server. See the [list of available images](https://websa.advancedhosting.com/slugs).
* `name` - (Required) Name for the Cloud Server.
* `datacenter` - (Required) Datacenter ID or Slug to start the Cloud Server in. See the [list of available datacenters](https://websa.advancedhosting.com/slugs).
* `product` - (**Deprecated**) Cloud Server Product ID or Slug that identifies the desired product type of the Cloud Server. See the [list of available products](https://websa.advancedhosting.com/slugs).
* `plan` - (Optional) Cloud Server Plan ID or Slug that identifies the desired plan type of the Cloud Server. Changing this resizes the existing server. See the [list of available products](https://websa.advancedhosting.com/slugs).
* `backups` - (Optional) Boolean to enable or disable backups. Defaults to false.
* `use_password` - (Optional) Boolean defining if password should be generated for the server and sent by email. Defaults to true.
* `ssh_keys` - (Optional) Array of SSH IDs or fingerprints to enable in
   the format `[12345, 7e:ac:a8:45:83:e3:58:f5:3a:9f:dd:16:63:dc:fb:1e]`. Fingerprints can be found in the 'SSH keys' section of the panel.
* `create_public_ip_address` - (Optional) Boolean defining if a new public IP address should be created for the server. This public IP address will become a primary IP address for the Cloud Server. Defaults to true.
* `private_cloud` (Optional) Boolean defining if instance should be created in private cloud
* `cluster_id` - (Optional, Required in case of `private_cloud=true`) The Cloud Server cluster ID
* `node_id` - (Optional, Required in case of `private_cloud=true`) The Cloud Server node ID
* `vcpu` - (Optional, Required in case of `private_cloud=true`) Required number of VCPUs for the Cloud Server  
* `ram` - (Optional, Required in case of `private_cloud=true`) Required RAM value for the Cloud Server 
* `disk` - (Optional, Required in case of `private_cloud=true`) Required disk size for the Cloud Server 

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
