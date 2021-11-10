<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for MSSQL Server (DBMS)

This MSSQL Server DBMS provider for Terraform allows you to manage databases and logins in an existing SQL Server installation.

The motivation for this provider is for use in automating ephemeral environments sharing an existing SQL Server instance.


## Usage Example

Requires Terraform 0.12.x and later, but 1.0 is recommended.

### Provider Configuration
> When using the provider with Terraform 0.13 and later, the recommended approach is to declare Provider versions in the root module Terraform configuration, using a `required_providers` block as per the following example. For previous versions, please continue to pin the version within the provider block.

```hcl
# We strongly recommend using the required_providers block to set the
# provider source and version being used
terraform {
  required_providers {
    mssqlserver = {
      source  = "github.com/ecabiac/mssqlserver"
      version = ">= 0.1.0"
    }
  }
}

# Configure the provider
provider "mssqlserver" {
  username  = "..."
  password  = "..."
  host  = "..."
  port = 1433
}
```

Further [usage documentation is available in the Examples directory](./examples/).

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x + (but 1.x is recommended)
* [Go](https://golang.org/doc/install) version 1.16.x (to build the provider plugin)

## Notes

  - Provider values must be known. See [Unknown values](https://www.terraform.io/docs/plugin/framework/providers.html#unknown-values)
  - Only use with trusted input. Not everything guards against SQL injection
  - Backup files must be available via the filesystem of the SQL Server Instance
