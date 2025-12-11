// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentcollection?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentregistry-post-agentcollections?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentcollection-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentcollection-update?view=graph-rest-beta
package graphBetaAgentsAgentCollection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentCollectionResourceModel represents the Terraform resource model for an agent collection.
type AgentCollectionResourceModel struct {
	// Required fields
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	OwnerIds    types.Set    `tfsdk:"owner_ids"`

	// Optional fields
	Description      types.String `tfsdk:"description"`
	ManagedBy        types.String `tfsdk:"managed_by"`
	OriginatingStore types.String `tfsdk:"originating_store"`

	// Computed fields (read-only)
	CreatedBy            types.String `tfsdk:"created_by"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`

	// Terraform-specific
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
