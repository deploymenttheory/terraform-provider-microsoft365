// REF: https://learn.microsoft.com/en-us/graph/api/resources/mobileappcatalogpackage?view=graph-rest-beta

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
