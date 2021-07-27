package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasaron/terraform-provider-stripe/stripe"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode,
		"debug",
		false,
		"set to true to run the provider with the debug support")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: stripe.Provider,
	}

	if debugMode {
		err := plugin.Debug(context.Background(), "local/lukasaron/stripe", opts)
		if err != nil {
			log.Fatal(err)
		}
	}

	plugin.Serve(opts)
}
