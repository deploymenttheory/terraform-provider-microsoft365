package graphBetaRoleScopeTag

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleScopeTagDataSourceModel defines the data source model
type RoleScopeTagDataSourceModel struct {
	FilterType  types.String        `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String        `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []RoleScopeTagModel `tfsdk:"items"`        // List of role scope tags that match the filters
	Timeouts    timeouts.Value      `tfsdk:"timeouts"`
}

// RoleScopeTagModel represents a single role scope tag
type RoleScopeTagModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
