// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentidentityblueprintprincipal?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/entra/agent-id/identity-platform/create-blueprint?tabs=microsoft-graph-api#create-a-service-principal-for-the-blueprint
// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities
package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentIdentityBlueprintServicePrincipalResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	AppId      types.String   `tfsdk:"app_id"`
	Tags       types.Set      `tfsdk:"tags"`
	HardDelete types.Bool     `tfsdk:"hard_delete"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}
