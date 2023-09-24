package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/myklst/terraform-provider-acme/acme"
	"github.com/myklst/terraform-provider-acme/acme/dnsplugin"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name acme

func main() {
	providerAddress := os.Getenv("PROVIDER_LOCAL_PATH")
	if providerAddress == "" {
		providerAddress = "registry.terraform.io/myklst/acme"
	}

	if len(os.Args) == 2 && os.Args[1] == dnsplugin.PluginArg {
		// Start the plugin here
		dnsplugin.Serve()
	} else {
		providerserver.Serve(context.Background(), acme.New, providerserver.ServeOpts{
			Address: providerAddress,
		})
	}
}
