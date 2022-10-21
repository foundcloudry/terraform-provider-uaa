package main

import (
	"flag"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		ProviderAddr: "registry.terraform.io/cloudfoundry-community/cloudfoundry",
		Debug:        debugMode,
	})
}
