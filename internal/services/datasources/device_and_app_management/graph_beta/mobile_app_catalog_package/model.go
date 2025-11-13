// REF: undocumented api endpoint
//
// IMPORTANT: Microsoft Graph API Limitations for mobileAppCatalogPackages endpoint
// This endpoint has significant limitations with OData query parameters that affect the model design.
//
// WORKING OData Parameters:
//   - odata_filter: Only works with startswith() function (e.g., "startswith(publisherDisplayName, 'Microsoft')")
//   - odata_top: Works for limiting results
//
// NOT WORKING/PROBLEMATIC OData Parameters:
//   - odata_skip: Causes 500 Internal Server errors and timeouts - DO NOT USE
//   - odata_select: Causes 500 Internal Server errors and timeouts - DO NOT USE
//   - odata_orderby: Returns no results when combined with odata_filter - DO NOT COMBINE
//   - odata_count: Returns no results when combined with odata_filter - DO NOT COMBINE
//   - odata_search: Not reliably supported by this endpoint - AVOID
//   - eq operator in filters: Not reliable, use startswith() instead
//
// RECOMMENDED USAGE:
//   - Use simple filter_type values: "all", "id", "product_name", "publisher_name"
//   - If using OData, only use odata_filter (with startswith()) and odata_top
//   - Avoid combining multiple OData parameters as they cause empty results or errors

package graphBetaMobileAppCatalogPackage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppCatalogPackageDataSourceModel defines the data source model
type MobileAppCatalogPackageDataSourceModel struct {
	FilterType   types.String                   `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String                   `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String                   `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32                    `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32                    `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String                   `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String                   `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	ODataCount   types.Bool                     `tfsdk:"odata_count"`   // OData count parameter
	ODataSearch  types.String                   `tfsdk:"odata_search"`  // OData search parameter
	ODataExpand  types.String                   `tfsdk:"odata_expand"`  // OData expand parameter
	Items        []MobileAppCatalogPackageModel `tfsdk:"items"`         // List of mobile app catalog packages that match the filters
	Timeouts     timeouts.Value                 `tfsdk:"timeouts"`
}

// MobileAppCatalogPackageModel represents a single mobile app catalog package
type MobileAppCatalogPackageModel struct {
	ID                       types.String   `tfsdk:"id"`
	ProductID                types.String   `tfsdk:"product_id"`
	ProductDisplayName       types.String   `tfsdk:"product_display_name"`
	PublisherDisplayName     types.String   `tfsdk:"publisher_display_name"`
	VersionDisplayName       types.String   `tfsdk:"version_display_name"`
	BranchDisplayName        types.String   `tfsdk:"branch_display_name"`
	ApplicableArchitectures  types.String   `tfsdk:"applicable_architectures"`
	Locales                  []types.String `tfsdk:"locales"`
	PackageAutoUpdateCapable types.Bool     `tfsdk:"package_auto_update_capable"`
}
