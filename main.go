package main

import (
	"github.com/esskar/terraform-provider-uniteddomains/uniteddomains"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: uniteddomains.Provider,
	})
}
