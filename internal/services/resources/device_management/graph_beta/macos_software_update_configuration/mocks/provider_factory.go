package mocks

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestUnitTestProtoV6ProviderFactories are used to instantiate a provider during
// unit tests. It accepts a TestVersion argument to specify which provider version
// should be instantiated. This factory function maps provider names to a
// function that returns a *tfprotov6.ProviderServer for unit testing.
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(provider.New("test")())
	},
}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance tests. It accepts a TestVersion argument to specify which provider version
// should be instantiated. This factory function maps provider names to a
// function that returns a *tfprotov6.ProviderServer for acceptance testing.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(provider.New(context.Background().Value("version").(string))())
	},
}
