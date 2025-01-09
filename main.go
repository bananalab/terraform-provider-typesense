package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-typesense/typesense"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name typesense

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "bananalab/terraform/typesense",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), typesense.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
