# Cloudbit Terraform Provider

This repository contains the Terraform provider implementation for the [Cloudbit](https://cloudbit.ch/) platform.

## Developing

In order to develop the provider, you need to tell Terraform to use the locally built provider instead of fetching it
from the registry. To do this, add the following to your `~/.terraformrc` file:

```hcl
provider_installation {
  # Use the given directory as the provider installation directory.
  # This disables the version and checksum verifications for this
  # provider and forces Terraform to look for the provider plugin
  # in the given directory.
  dev_overrides {
    "cloudbit-ch/cloudbit" = "path to local `terraform-provider-cloudbit` directory"
  }

  # For all other providers, use the default behavior.
  direct {}
}
```

Please see the
[Terraform Documentation](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers)
for more details about the `dev_overrides` section.

Once you have configured your `~/.terraformrc`, you must build the provider every time you change your code using
`go build .`. This generates the `terraform-provider-cloudbit` binary which terraform can then use as a provider.