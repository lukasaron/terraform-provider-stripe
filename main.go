package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasaron/terraform-provider-stripe/stripe"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: stripe.Provider})
}
