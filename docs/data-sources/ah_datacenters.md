# AH Datacenters Data Source

Get information about AdvancedHosting Datacenters.

## Example Usage

Get the Datacenters by ID:

```hcl
data "ah_datacenters" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a list of Datacenters in Netherlands sorted by slug:

```hcl
data "ah_datacenters" "example" {
  filter {
    key = "region_country_code"
    values = ["NL"]
  }
  sort {
    key = "slug"
    direction = "asc"
  }
}
```

## Argument Reference

* `filter`: (Optional) Filter the results by specified key and value. The structure of the block is documented below.
* `sort` - (Optional) Sort the results by specified key and direction. The structure of the block is documented below.

---

The `filter` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `full_name`, `region_id`, `region_name`, `region_country_code`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `slug`, `full_name`, `region_id`, `region_name`, `region_country_code`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `datacenters` - A list of Datacenters that satisfy the search criteria.
    * `id` -  ID of the Datacenter.
    * `name` - Datacenter name.
    * `slug` - Datacenter slug.
    * `full_name` - Datacenter full name.
    * `region_id` -  Datacenter region ID.
    * `region_name` - Datacenter region name.
    * `region_country_code` - Datacenter region country code.
