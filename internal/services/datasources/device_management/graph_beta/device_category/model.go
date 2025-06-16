package graphBetaDeviceCategory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceCategoryDataSourceModel defines the data source model
type DeviceCategoryDataSourceModel struct {
	FilterType  types.String          `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String          `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []DeviceCategoryModel `tfsdk:"items"`        // List of device categories that match the filters
	Timeouts    timeouts.Value        `tfsdk:"timeouts"`
}

// DeviceCategoryModel represents a single device category
type DeviceCategoryModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
