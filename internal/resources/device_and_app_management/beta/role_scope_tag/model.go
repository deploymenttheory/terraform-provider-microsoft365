package graphBetaRoleScopeTag

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleScopeTagResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-rolescopetag?view=graph-rest-beta
type RoleScopeTagResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	DisplayName types.String   `tfsdk:"display_name"`
	Description types.String   `tfsdk:"description"`
	IsBuiltIn   types.Bool     `tfsdk:"is_built_in"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}
