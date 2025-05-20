// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta

package graphBetaMobileApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppDataSourceModel defines the data source model
type MobileAppDataSourceModel struct {
	FilterType    types.String     `tfsdk:"filter_type"`     // Required field to specify how to filter
	FilterValue   types.String     `tfsdk:"filter_value"`    // Value to filter by (not used for "all" or "odata")
	ODataFilter   types.String     `tfsdk:"odata_filter"`    // OData filter parameter
	ODataTop      types.Int32      `tfsdk:"odata_top"`       // OData top parameter for limiting results
	ODataSkip     types.Int32      `tfsdk:"odata_skip"`      // OData skip parameter for pagination
	ODataSelect   types.String     `tfsdk:"odata_select"`    // OData select parameter for field selection
	ODataOrderBy  types.String     `tfsdk:"odata_orderby"`   // OData orderby parameter for sorting
	AppTypeFilter types.String     `tfsdk:"app_type_filter"` // Optional filter to filter by odata mobile app type
	Items         []MobileAppModel `tfsdk:"items"`           // List of mobile apps that match the filters
	Timeouts      timeouts.Value   `tfsdk:"timeouts"`
}

// MobileAppModel represents a single mobile app with common fields
type MobileAppModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	Description           types.String   `tfsdk:"description"`
	Publisher             types.String   `tfsdk:"publisher"`
	Developer             types.String   `tfsdk:"developer"`
	Owner                 types.String   `tfsdk:"owner"`
	Notes                 types.String   `tfsdk:"notes"`
	CreatedDateTime       types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String   `tfsdk:"last_modified_date_time"`
	InformationUrl        types.String   `tfsdk:"information_url"`
	PrivacyInformationUrl types.String   `tfsdk:"privacy_information_url"`
	PublishingState       types.String   `tfsdk:"publishing_state"`
	IsAssigned            types.Bool     `tfsdk:"is_assigned"`
	IsFeatured            types.Bool     `tfsdk:"is_featured"`
	UploadState           types.Int64    `tfsdk:"upload_state"`
	DependentAppCount     types.Int64    `tfsdk:"dependent_app_count"`
	SupersededAppCount    types.Int64    `tfsdk:"superseded_app_count"`
	SupersedingAppCount   types.Int64    `tfsdk:"superseding_app_count"`
	RoleScopeTagIds       []types.String `tfsdk:"role_scope_tag_ids"`
	Categories            []types.String `tfsdk:"categories"`
}
