package graphBetaAssignmentFilter

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AssignmentFilterDataSourceModel defines the data source model
type AssignmentFilterDataSourceModel struct {
	FilterType  types.String            `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String            `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []AssignmentFilterModel `tfsdk:"items"`        // List of assignment filters that match the filters
	Timeouts    timeouts.Value          `tfsdk:"timeouts"`
}

// AssignmentFilterModel represents a single assignment filter
type AssignmentFilterModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}
