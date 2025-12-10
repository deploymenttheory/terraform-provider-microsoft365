// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentidentity?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities
package graphBetaAgentIdentity

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentIdentityResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	DisplayName               types.String   `tfsdk:"display_name"`
	AgentIdentityBlueprintId  types.String   `tfsdk:"agent_identity_blueprint_id"`
	AccountEnabled            types.Bool     `tfsdk:"account_enabled"`
	CreatedByAppId            types.String   `tfsdk:"created_by_app_id"`
	CreatedDateTime           types.String   `tfsdk:"created_date_time"`
	DisabledByMicrosoftStatus types.String   `tfsdk:"disabled_by_microsoft_status"`
	ServicePrincipalType      types.String   `tfsdk:"service_principal_type"`
	Tags                      types.Set      `tfsdk:"tags"`
	SponsorIds                types.Set      `tfsdk:"sponsor_ids"`
	OwnerIds                  types.Set      `tfsdk:"owner_ids"`
	HardDelete                types.Bool     `tfsdk:"hard_delete"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
}
