package graphBetaLinuxPlatformScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinuxPlatformScriptDataSourceModel defines the data source model
type LinuxPlatformScriptDataSourceModel struct {
	FilterType  types.String               `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String               `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []LinuxPlatformScriptModel `tfsdk:"items"`        // List of Linux platform scripts that match the filters
	Timeouts    timeouts.Value             `tfsdk:"timeouts"`
}

// LinuxPlatformScriptModel represents a single Linux platform script
type LinuxPlatformScriptModel struct {
	ID           types.String   `tfsdk:"id"`
	DisplayName  types.String   `tfsdk:"display_name"`
	Description  types.String   `tfsdk:"description"`
	Technologies []types.String `tfsdk:"technologies"`
}
