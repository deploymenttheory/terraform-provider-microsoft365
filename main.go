package main

import (
	"context"
	"flag"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	version string = "dev"
	commit  string = "unknown"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name microsoft365

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	tflog.Info(ctx, "Starting Microsoft365 provider", map[string]interface{}{
		"version": version,
		"commit":  commit,
	})

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/deploymenttheory/microsoft365",
		Debug:   debug,
	}

	tflog.Debug(ctx, "Provider serve options", map[string]interface{}{
		"address": opts.Address,
		"debug":   opts.Debug,
	})

	err := providerserver.Serve(ctx, provider.NewMicrosoft365Provider(version), opts)

	if err != nil {
		tflog.Error(ctx, "Provider server error", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
