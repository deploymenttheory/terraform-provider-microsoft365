package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	providercore "github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider/core"
)

// M365Provider combines the shared provider configuration with all production registrations.
type M365Provider struct {
	*providercore.M365Provider
}

// NewMicrosoft365Provider creates the complete production provider.
func NewMicrosoft365Provider(version string, unitTestMode ...bool) func() provider.Provider {
	return func() provider.Provider {
		return &M365Provider{M365Provider: providercore.NewMicrosoft365Provider(version, unitTestMode...)().(*providercore.M365Provider)}
	}
}
