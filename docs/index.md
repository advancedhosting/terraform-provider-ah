---
layout: "ah"
page_title: "Provider: Advahced Hosting"
sidebar_current: "docs-ah-index"
description: |-
    The Advahced Hosting provider is used to manage Advahced Hosting resources. The provider needs to be configured with the proper credentials before it can be used.

---

# Advanced Hosting Provider

The Advanced Hosting terraform provider is used to interact with Advanced Hosting (AH) resources. Can be used to create, modify, and delete Cloud Servers and other resources. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the Advanced Hosting provider
provider "ah" {
  token = "auth_token_here"
}

// Create a new instance
resource "ah_cloud_server" "default" {
  ...
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) Security token used for authentication in Advanced Hosting. This can also be specified using environment variable `AH_TOKEN`.
* `api_endpoint` - (Optional) Specify which API endpoint to use, can be used to override the default API Endpoint. This can also be specified using environment variable `AH_API_ENDPOINT`. 

