package graphBetaWindowsUpdateRing

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsUpdateRingDataSourceModel defines the data source model
type WindowsUpdateRingDataSourceModel struct {
	FilterType  types.String             `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String             `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []WindowsUpdateRingModel `tfsdk:"items"`        // List of Windows Update Rings that match the filters
	Timeouts    timeouts.Value           `tfsdk:"timeouts"`
}

// WindowsUpdateRingModel represents a single Windows Update Ring
type WindowsUpdateRingModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
