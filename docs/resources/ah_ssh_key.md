---
layout: "ah"
page_title: "Advahced Hosting: ah_ssh_key"
sidebar_current: "docs-resource-ah-ssh-key"
description: |-
    Manages Advanced Hosting SSH keys.
---

# ah_ssh_key

Provides an Advanced Hosting SSH key resource.

## Example Usage

```hcl

resource "ah_ssh_key" "example" {
  name = "SSH key Name"
  public_key = file("~/.ssh/id_rsa.pub")
}

```

## Argument Reference

The following arguments are supported:
* `name` - (Optional) SSH key name
* `public_key` - (Required) Public key. If this is a file, it can be read using the file interpolation function.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of the SSH key.
* `fingerprint` - Fingerprint of the SSH key.
* `created_at` - Creation datetime of the SSH key.
