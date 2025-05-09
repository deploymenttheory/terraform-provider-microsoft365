// REF: https://learn.microsoft.com/en-us/graph/api/resources/browsersite?view=graph-rest-beta
package graphBetaBrowserSite

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BrowserSiteDataSourceModel defines the data source model
type BrowserSiteDataSourceModel struct {
	FilterType                  types.String               `tfsdk:"filter_type"`                     // Required field to specify how to filter
	FilterValue                 types.String               `tfsdk:"filter_value"`                    // Value to filter by (not used for "all")
	BrowserSiteListAssignmentID types.String               `tfsdk:"browser_site_list_assignment_id"` // Required BrowserSiteList ID
	Items                       []BrowserSiteResourceModel `tfsdk:"items"`                           // List of items that match the filters
	Timeouts                    timeouts.Value             `tfsdk:"timeouts"`
}

// BrowserSiteResourceModel represents a simplified browser site for the data source
type BrowserSiteResourceModel struct {
	ID                          types.String `tfsdk:"id"`                              // The unique identifier of the browser site
	BrowserSiteListAssignmentID types.String `tfsdk:"browser_site_list_assignment_id"` // The browser site list ID
	WebUrl                      types.String `tfsdk:"web_url"`                         // The URL of the site
}
