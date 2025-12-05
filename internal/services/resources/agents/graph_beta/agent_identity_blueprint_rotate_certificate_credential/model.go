// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addkey?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-removekey?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/keycredential?view=graph-rest-beta
package graphBetaAgentIdentityBlueprintKeyCredential

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentIdentityBlueprintKeyCredentialResourceModel describes the resource data model.
type AgentIdentityBlueprintKeyCredentialResourceModel struct {
	// Required inputs
	BlueprintID types.String `tfsdk:"blueprint_id"`
	Key         types.String `tfsdk:"key"`
	KeyType     types.String `tfsdk:"type"`
	Usage       types.String `tfsdk:"usage"`
	Proof       types.String `tfsdk:"proof"`

	// Optional inputs
	DisplayName         types.String `tfsdk:"display_name"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	PasswordSecretText  types.String `tfsdk:"password_secret_text"`

	// Computed outputs (from API response)
	KeyID types.String `tfsdk:"key_id"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
