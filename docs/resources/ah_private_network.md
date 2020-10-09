# AH Private Network Resource

Provides an Advanced Hosting Private Network resource to represent a Private Network that can be connect to a Cloud Server.


## Example Usage

```hcl

resource "ah_private_network" "example" {
  ip_range = "10.0.0.0/24"
  name = "New Private Network"
}

```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) Private Network IP range in CIDR format.
* `name` - (Required) Name of the Private Network.

---

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - ID of Private Network.
* `state` - Current state of the Private Network.
* `created_at` - Creation datetime of the Private Network.