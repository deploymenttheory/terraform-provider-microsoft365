package graphBetaReuseablePolicySettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReuseablePolicySettingsDataSourceModel defines the data source model
type ReuseablePolicySettingsDataSourceModel struct {
	FilterType  types.String                  `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                  `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []ReuseablePolicySettingModel `tfsdk:"items"`        // List of reusable policy settings that match the filters
	Timeouts    timeouts.Value                `tfsdk:"timeouts"`
}

// ReuseablePolicySettingModel represents a single reusable policy setting
type ReuseablePolicySettingModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
