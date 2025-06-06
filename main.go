package main

import (
	"context"
	"flag"
	"github.com/lukasaron/terraform-provider-stripe/internal/provider"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with debugging support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/lukasaron/stripe",
		Debug:   debugMode,
	}

	if err := providerserver.Serve(context.Background(), provider.New(), opts); err != nil {
		log.Fatal(err)
	}
}
