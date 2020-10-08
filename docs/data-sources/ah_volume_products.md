# AH Volume Products Data Source

Get information about AdvancedHosting Volume Products available for volume creation.

## Example Usage

Get the Volume Product by ID:

```hcl
data "ah_volume_products" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of Volume available in AMS1 datacenter, sorted by maximum volume size, desc:

```hcl
data "ah_volume_products" "example" {
  filter {
    key = "datacenter_slug"
    values = ['ams1']
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

* `products` - A list of Products that satisfy the search criteria.
  * `id` - ID of the Product.
  * `name` - Name of the Product.
  * `slug` - Slug of the Product.
  * `price` - Price of the Product (per GB/month).
  * `currency` - Currency for the price.
  * `min_size` - Minimum size available for Volume creation in GB.
  * `max_size` - Maximum size available for Volume creation in GB.
  * `datacenters`- List of Cloud Server Datacenters the Volume can be attached to. The structure of the block is documented below.

---

The `datacenters` block contains:
* `id` - ID of the Datacenter.
* `name` - Datacenter name.
* `slug` - Datacenter slug.
* `full_name` - Datacenter full name.