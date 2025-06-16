package graphBetaWindowsRemediationScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsRemediationScriptDataSourceModel defines the data source model
type WindowsRemediationScriptDataSourceModel struct {
	FilterType  types.String                    `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                    `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []WindowsRemediationScriptModel `tfsdk:"items"`        // List of Windows Remediation Scripts that match the filters
	Timeouts    timeouts.Value                  `tfsdk:"timeouts"`
}

// WindowsRemediationScriptModel represents a single Windows Remediation Script
type WindowsRemediationScriptModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
