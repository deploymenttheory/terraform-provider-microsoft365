// REF: https://learn.microsoft.com/en-us/graph/api/agentcollection-list-members?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentcollection-post-members?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentcollection-delete-members?view=graph-rest-beta
package graphBetaAgentsAgentCollectionAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentCollectionAssignmentResourceModel represents the Terraform resource model for agent collection membership.
type AgentCollectionAssignmentResourceModel struct {
	// Required fields
	ID                types.String `tfsdk:"id"`
	AgentInstanceID   types.String `tfsdk:"agent_instance_id"`
	AgentCollectionID types.String `tfsdk:"agent_collection_id"`

	// Terraform-specific
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
