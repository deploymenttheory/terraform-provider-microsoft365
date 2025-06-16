// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-policyset-deviceandappmanagementassignmentfilter?view=graph-rest-beta
package graphBetaAssignmentFilter

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AssignmentFilterResourceModel struct {
	ID                             types.String   `tfsdk:"id"`
	DisplayName                    types.String   `tfsdk:"display_name"`
	Description                    types.String   `tfsdk:"description"`
	Platform                       types.String   `tfsdk:"platform"`
	Rule                           types.String   `tfsdk:"rule"`
	AssignmentFilterManagementType types.String   `tfsdk:"assignment_filter_management_type"`
	CreatedDateTime                types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime           types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTags                  types.Set      `tfsdk:"role_scope_tags"`
	Timeouts                       timeouts.Value `tfsdk:"timeouts"`
}
