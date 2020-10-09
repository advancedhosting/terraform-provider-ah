# Terraform Provider for the AdvancedHosting Cloud

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

## Building The Provider

Clone repository and run:

```
make install
```

To learn more about using a local build of a provider, you can look at the [documentation on writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html#invoking-the-provider) and the [documentation on how Terraform plugin discovery works](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery)

Using the provider
----------------------
See the documentation to get started using the AdvancedHosting provider.