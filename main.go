package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasaron/terraform-provider-stripe/stripe"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: stripe.Provider,
	}
	// debug mode is on for now
	//err := plugin.Debug(context.Background(), "local/lukasaron/stripe", opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return

	// TODO uncomment when debugging is not needed.
	plugin.Serve(opts)
}
