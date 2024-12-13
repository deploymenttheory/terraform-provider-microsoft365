package graphBetaRoleScopeTag

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleScopeTagResourceModel holds the configuration for a Role Scope Tag.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-rolescopetag?view=graph-rest-beta
type RoleScopeTagResourceModel struct {
	ID          types.String      `tfsdk:"id"`
	DisplayName types.String      `tfsdk:"display_name"`
	Description types.String      `tfsdk:"description"`
	IsBuiltIn   types.Bool        `tfsdk:"is_built_in"`
	Assignments []AssignmentModel `tfsdk:"assignments"`
	Timeouts    timeouts.Value    `tfsdk:"timeouts"`
}

// AssignmentModel represents a single assignment for a role scope tag
type AssignmentModel struct {
	GroupID types.String `tfsdk:"group_id"`
}
