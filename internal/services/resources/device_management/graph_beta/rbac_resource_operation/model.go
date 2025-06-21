// resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-RBACResourceOperation?view=graph-rest-beta
package graphBetaRBACResourceOperation

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RBACResourceOperationResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	Resource                  types.String   `tfsdk:"resource"`
	ResourceName              types.String   `tfsdk:"resource_name"`
	ActionName                types.String   `tfsdk:"action_name"`
	Description               types.String   `tfsdk:"description"`
	EnabledForScopeValidation types.Bool     `tfsdk:"enabled_for_scope_validation"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
}
