// REF: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-list-appliesto?view=graph-rest-beta
package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type B2bManagementPolicyAssignmentResourceModel struct {
	ID                    types.String   `tfsdk:"id"` // "{b2b_management_policy_id}/{directory_object_id}"; Graph has no assignment ID
	B2bManagementPolicyID types.String   `tfsdk:"b2b_management_policy_id"`
	DirectoryObjectID     types.String   `tfsdk:"directory_object_id"`
	DirectoryObjectType   types.String   `tfsdk:"directory_object_type"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
