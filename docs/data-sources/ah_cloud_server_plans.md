# AH Cloud Server Plans Data Source

Get information about AdvancedHosting Cloud Server Plans available for server creation.

## Example Usage

Get the Cloud Server Plan by ID:

```hcl
data "ah_cloud_server_plans" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of Plans that have 2 or 3 vCPUs, sorted by RAM, desc:

```hcl
data "ah_cloud_server_plans" "example" {
  filter {
    key = "vcpu"
    values = [2, 3]
  }
  sort {
    key = "ram"
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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `price`, `currency`, `vcpu`, `ram`, `disk`, `available_on_trial`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `price`, `currency`, `vcpu`, `ram`, `disk`, `available_on_trial`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `plans` - A list of Plans that satisfy the search criteria.
  * `id` - ID of the Plan.
  * `name` - Name of the Plan.
  * `slug` - Slug of the Plan.
  * `price` - Monthly price of the Plan.
  * `currency` - Currency for the price.
  * `vcpu` - Number of vCPUs.
  * `ram` - RAM in GiB.
  * `disk` - Disk size in GB.
  * `available_on_trial` - Boolean flag indicating whether the Plan is available on trial.
