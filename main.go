package main

import (
	"github.com/esskar/terraform-provider-uniteddomains/uniteddomains"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return uniteddomains.Provider()
		},
	})
}
