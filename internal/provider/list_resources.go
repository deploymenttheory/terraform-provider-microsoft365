package provider

import (
	"context"

	graphBetaSettingsCatalogConfigurationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ProviderWithListResources = &M365Provider{}

// ListResources defines the list resources implemented in the provider.
func (p *M365Provider) ListResources(_ context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		// Graph Beta - Device Management list resources
		graphBetaSettingsCatalogConfigurationPolicy.NewSettingsCatalogListResource,

		// Add microsoft 365 provider list resources here
	}
}
