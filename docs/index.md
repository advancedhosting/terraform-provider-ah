---
layout: "ah"
page_title: "Provider: Advahced Hosting"
sidebar_current: "docs-ah-index"
description: |-
    The Advahced Hosting provider is used to manage Advahced Hosting resources. The provider needs to be configured with the proper credentials before it can be used.

---

# AdvancedHosting Provider

The AdvancedHosting terraform provider is used to interact with AdvancedHosting (AH) resources. Can be used to create, modify, and delete Cloud Servers and other resources. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the AdvancedHosting provider
provider "ah" {
  access_token = "auth_token_here"
}

// Create a new instance
resource "ah_cloud_server" "default" {
  ...
}
```

## Argument Reference

The following arguments are supported:

* `access_token` - (Required) Security token used for authentication in AdvancedHosting. This can also be specified using the environment variable `AH_ACCESS_TOKEN`.
* `endpoint` - (Optional) Specify which API endpoint to use, can be used to override the default API Endpoint. This can also be specified using the environment variable `AH_API_ENDPOINT`. 

