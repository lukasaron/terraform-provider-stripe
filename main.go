package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasaron/terraform-provider-stripe/stripe"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: stripe.Provider,
	}

	plugin.Serve(opts)
}
