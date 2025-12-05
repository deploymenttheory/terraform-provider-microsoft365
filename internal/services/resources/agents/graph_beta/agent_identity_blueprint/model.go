// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentidentityblueprint?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/entra/agent-id/identity-platform/create-blueprint?tabs=microsoft-graph-api
// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities
package graphBetaApplicationsAgentIdentityBlueprint

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentIdentityBlueprintResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	AppId          types.String   `tfsdk:"app_id"`
	DisplayName    types.String   `tfsdk:"display_name"`
	Description    types.String   `tfsdk:"description"`
	SignInAudience types.String   `tfsdk:"sign_in_audience"`
	Tags           types.Set      `tfsdk:"tags"`
	SponsorUserIds types.Set      `tfsdk:"sponsor_user_ids"`
	OwnerUserIds   types.Set      `tfsdk:"owner_user_ids"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}
