// REF: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-list-appliesto?view=graph-rest-beta
package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// B2bManagementPolicyAssignmentResourceModel represents the schema for the B2B management policy assignment resource.
// The Microsoft Graph API does not return an ID for this assignment; ID is a provider-constructed composite of
// b2b_management_policy_id and directory_object_id, used to track the assignment in Terraform state.
type B2bManagementPolicyAssignmentResourceModel struct {
	ID                    types.String   `tfsdk:"id"` // synthetic: "{b2b_management_policy_id}/{directory_object_id}"
	B2bManagementPolicyID types.String   `tfsdk:"b2b_management_policy_id"`
	DirectoryObjectID     types.String   `tfsdk:"directory_object_id"`
	DirectoryObjectType   types.String   `tfsdk:"directory_object_type"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
