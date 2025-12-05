// REF: https://learn.microsoft.com/en-us/entra/agent-id/identity-platform/create-blueprint?tabs=microsoft-graph-api
package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentIdentityBlueprintIdentifierUriResourceModel describes the resource data model.
type AgentIdentityBlueprintIdentifierUriResourceModel struct {
	// Required inputs
	BlueprintID   types.String `tfsdk:"blueprint_id"`
	IdentifierUri types.String `tfsdk:"identifier_uri"`

	// Optional: OAuth2 Permission Scope configuration
	Scope *OAuth2PermissionScopeModel `tfsdk:"scope"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// OAuth2PermissionScopeModel describes the OAuth2 permission scope configuration.
type OAuth2PermissionScopeModel struct {
	ID                      types.String `tfsdk:"id"`
	AdminConsentDescription types.String `tfsdk:"admin_consent_description"`
	AdminConsentDisplayName types.String `tfsdk:"admin_consent_display_name"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
	Type                    types.String `tfsdk:"type"`
	Value                   types.String `tfsdk:"value"`
}
