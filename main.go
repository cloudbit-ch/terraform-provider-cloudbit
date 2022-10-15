package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/cloudbit-ch/terraform-provider-cloudbit/cloudbit"
)

var (
	version         = "dev"
	defaultEndpoint = "https://api.cloudbit.ch/"
)

func main() {
	debug := false
	showVersion := false

	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&showVersion, "version", false, "show version and quit")
	flag.Parse()

	if showVersion {
		fmt.Println("terraform-provider-cloudbit", version)
		return
	}

	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/cloudbit-ch/cloudbit",
		Debug:           debug,
		ProtocolVersion: 6,
	}

	factory := func() tfsdk.Provider {
		return cloudbit.New(
			cloudbit.WithVersion(version),
			cloudbit.WithDefaultEndpoint(defaultEndpoint),
		)
	}

	err := providerserver.Serve(context.Background(), factory, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
