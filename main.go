package main

import (
	_ "github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/ecabiac/terraform-provider-mssqlserver/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
