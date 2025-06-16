// REF: https://learn.microsoft.com/en-us/graph/api/resources/browsersitelist?view=graph-rest-beta
package graphBetaBrowserSiteList

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BrowserSiteListDataSourceModel defines the data source model
type BrowserSiteListDataSourceModel struct {
	FilterType  types.String                   `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                   `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []BrowserSiteListResourceModel `tfsdk:"items"`        // List of items that match the filters
	Timeouts    timeouts.Value                 `tfsdk:"timeouts"`
}

// BrowserSiteListResourceModel represents a single catalog item
type BrowserSiteListResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}
