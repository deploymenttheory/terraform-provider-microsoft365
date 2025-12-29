// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-WindowsFeatureUpdateprofile?view=graph-rest-beta
package graphBetaWindowsFeatureUpdatePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsFeatureUpdateProfileDataSourceModel defines the data source model
type WindowsFeatureUpdateProfileDataSourceModel struct {
	ID              types.String   `tfsdk:"id"`
	DisplayName     types.String   `tfsdk:"display_name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.Set      `tfsdk:"role_scope_tag_ids"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}
