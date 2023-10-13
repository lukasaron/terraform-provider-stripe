package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/lukasaron/terraform-provider-stripe/stripe"
)

func main() {
	var debugMode bool
	var providerAddress string

	flag.BoolVar(&debugMode,
		"debug",
		false,
		"set to true to run the provider with the debug support")
	flag.StringVar(&providerAddress,
		"providerAddress",
		"lukasaron/stripe",
		"set the provider address to run the provider with the debug support")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: stripe.Provider,
		Debug:        debugMode,
		ProviderAddr: providerAddress,
	}

	plugin.Serve(opts)
}
