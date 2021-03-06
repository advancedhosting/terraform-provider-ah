
# AH Cloud Servers Data Source

Get information about AdvancedHosting Cloud Servers.

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

Get a list of active Cloud Servers from AMS1 and ASH1 datacenter, sorted by creation date in descending order:

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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `state`,  `vcpu`, `ram`, `disk`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `state`, `created_at`, `vcpu`, `ram`, `disk`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `cloud_servers` - A list of Cloud Servers that satisfy the search criteria.
  * `id` -  ID of the the Cloud Server.
  * `name` - Name of the Cloud Server.
  * `datacenter` - Datacenter slug of the Cloud Server.
  * `product` - Product slug that indentifies the product type of the Cloud Server.
  * `state` - Current state of the Cloud Server.
  * `vcpu` - Number of vCPUs on the Cloud Server.
  * `ram` - RAM of the server in MiB.
  * `disk` - Disk size of the server in GB.
  * `created_at` - Creation timestamp of the Cloud Server.
  * `image` - The Cloud Server Image ID or Snapshot / Auto Backup ID the server was created from.
  * `backups` - Boolean indicating whether backups are enabled for the Cloud Server.
  * `use_password` - Boolean indicating whether the Cloud Server was created with a password generated.
  * `ips` - Array of Public and Anycast IP addresses assigned to the the Cloud Server. The structure of the block is documented below.
  * `volumes` - Array of Volume IDs attached to the server.
  * `private_networks` - Array of Private Networks connected to the server. The structure of the block is documented below.
  
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
