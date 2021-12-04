# AH Volume Plans Data Source

Get information about AdvancedHosting Volume Plans available for volume creation.

## Example Usage

Get the Volume Plan by ID:

```hcl
data "ah_volume_plans" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of Volume Plans available in AMS1 datacenter, sorted by maximum volume size, desc:

```hcl
data "ah_volume_plans" "example" {
  filter {
    key = "datacenter_slug"
    values = ["ams1"]
  }
  sort {
    key = "max_size"
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
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `price`, `currency`, `min_size`, `max_size`, `datacenter_id`, `datacenter_name`, `datacenter_slug`, `datacenter_full_name`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `price`, `currency`, `min_size`, `max_size`, `datacenter_id`, `datacenter_name`, `datacenter_slug`, `datacenter_full_name`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `plans` - A list of Products that satisfy the search criteria.
  * `id` - ID of the Plan.
  * `name` - Name of the Plan.
  * `slug` - Slug of the Plan.
  * `price` - Price of the Plan (per GB/month).
  * `currency` - Currency for the price.
  * `min_size` - Minimum size available for Volume creation in GB.
  * `max_size` - Maximum size available for Volume creation in GB.
  * `datacenter_id`- ID of the Datacenter. 
  * `datacenter_slig`- Slug of the Datacenter. 