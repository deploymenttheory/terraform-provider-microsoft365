package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsQualityUpdateExpeditePolicyDataSourceModel defines the data source model
type WindowsQualityUpdateExpeditePolicyDataSourceModel struct {
	FilterType  types.String                              `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                              `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []WindowsQualityUpdateExpeditePolicyModel `tfsdk:"items"`        // List of Windows Quality Update Expedite Policies that match the filters
	Timeouts    timeouts.Value                            `tfsdk:"timeouts"`
}

// WindowsQualityUpdateExpeditePolicyModel represents a single Windows Quality Update Expedite Policy
type WindowsQualityUpdateExpeditePolicyModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
