// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addpassword?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-removepassword?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/entra/agent-id/identity-platform/create-blueprint?tabs=microsoft-graph-api
package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentIdentityBlueprintPasswordCredentialResourceModel describes the resource data model.
type AgentIdentityBlueprintPasswordCredentialResourceModel struct {
	// Required inputs
	BlueprintID types.String `tfsdk:"blueprint_id"`
	DisplayName types.String `tfsdk:"display_name"`

	// Optional inputs
	StartDateTime types.String `tfsdk:"start_date_time"`
	EndDateTime   types.String `tfsdk:"end_date_time"`

	// Computed outputs (from API response)
	KeyID               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	Hint                types.String `tfsdk:"hint"`
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
