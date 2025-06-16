// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsupdatecatalogitem?view=graph-rest-beta
package graphBetaWindowsUpdateCatalogItem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsUpdateCatalogItemDataSourceModel defines the data source model
type WindowsUpdateCatalogItemDataSourceModel struct {
	FilterType  types.String                    `tfsdk:"filter_type"`  // Required field to specify how to filter
	FilterValue types.String                    `tfsdk:"filter_value"` // Value to filter by (not used for "all")
	Items       []WindowsUpdateCatalogItemModel `tfsdk:"items"`        // List of catalog items that match the filters
	Timeouts    timeouts.Value                  `tfsdk:"timeouts"`
}

// WindowsUpdateCatalogItemModel represents a single catalog item
type WindowsUpdateCatalogItemModel struct {
	ID               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	ReleaseDateTime  types.String `tfsdk:"release_date_time"`
	EndOfSupportDate types.String `tfsdk:"end_of_support_date"`
}
