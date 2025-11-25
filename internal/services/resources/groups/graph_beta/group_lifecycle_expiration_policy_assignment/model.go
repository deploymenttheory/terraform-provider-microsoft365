// REF: https://learn.microsoft.com/en-us/graph/api/grouplifecyclepolicy-addgroup?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/grouplifecyclepolicy-removegroup?view=graph-rest-beta
package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupLifecycleExpirationPolicyAssignmentResourceModel represents the Terraform resource model
// for managing group assignments to the tenant's lifecycle policy.
// This resource is only applicable when the policy's managedGroupTypes is set to "Selected".
// Since there is only one global lifecycle policy per tenant, the resource ID is simply the group_id.
type GroupLifecycleExpirationPolicyAssignmentResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	GroupID  types.String   `tfsdk:"group_id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
