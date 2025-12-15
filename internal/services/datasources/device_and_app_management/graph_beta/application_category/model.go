package graphBetaApplicationCategory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationCategoryDataSourceModel defines the data source model
type ApplicationCategoryDataSourceModel struct {
	FilterType  types.String               `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String               `tfsdk:"filter_value"` // Value to filter by (not used for "all" or "odata")
	ODataFilter types.String               `tfsdk:"odata_filter"` // OData filter query (only used when filter_type is "odata")
	ODataTop    types.Int32                `tfsdk:"odata_top"`    // OData $top parameter (only used when filter_type is "odata")
	Items       []ApplicationCategoryModel `tfsdk:"items"`        // List of application categories that match the filters
	Timeouts    timeouts.Value             `tfsdk:"timeouts"`
}

// ApplicationCategoryModel represents a single application category
type ApplicationCategoryModel struct {
	ID                   types.String `tfsdk:"id"`
	DisplayName          types.String `tfsdk:"display_name"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
}
