---
layout: "ah"
page_title: "Advahced Hosting: ah_ssh_keys"
sidebar_current: "docs-data-source-ah-ssh-keys"
description: |-
  Get information about Advanced Hosting SSH keys.
---

# ah_ssh_keys

Get information about Advanced Hosting SSH keys.

## Example Usage

Get a SSH key by ID:

```hcl
data "ah_ssh_keys" "example" {
  filter {
    key = "id"
    values = ["123"]
  }
}
```

Get a SSH key by a Fingerprint:

```hcl
data "ah_ssh_keys" "example" {
  filter {
    key = "fingerprint"
    values = ["7e:ac:a8:45:83:e3:58:f5:3a:9f:dd:16:63:dc:fb:1e"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter`: (Optional) Filter the results by specified key and value. The structure of the block is documented below.
* `sort` - (Optional) Sort the results by specified key and direction. The structure of the block is documented below.

---

The `filter` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `fingerprint`, `created_at`
* `values` - (Required) A list of values to match against the `key` field.

The `sort` block supports:
* `key` - (Required) Filter the results by specified key. Can be one of: `id`, `name`, `fingerprint`, `created_at`
* `direction` - (Optional) Sort direction of the results. Can be one of: `asc`, `desc`. Default option is `desc`.

---

## Attributes Reference

The following attributes are exported:

* `ssh_keys` - A list of Private Networks that satisfy the search criteria.
    * `id` - ID of the SSH key.
    * `name` - SSH key name.
    * `public_key` - Public key.
    * `fingerprint` - Fingerprint of the SSH key.
    * `created_at` - Creation datetime of the SSH key.