// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta

package graphBetaMacOSPKGApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSPKGAppDataSourceModel defines the data source model
type MacOSPKGAppDataSourceModel struct {
	FilterType  types.String       `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String       `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []MacOSPKGAppModel `tfsdk:"items"`        // List of macOS PKG apps that match the filters
	Timeouts    timeouts.Value     `tfsdk:"timeouts"`
}

// MacOSPKGAppModel represents a single macOS PKG app
type MacOSPKGAppModel struct {
	ID              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	RoleScopeTagIds types.Set    `tfsdk:"role_scope_tag_ids"`
}
