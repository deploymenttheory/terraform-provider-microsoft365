// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta

package graphBetaMobileApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppDataSourceModel represents the Terraform data source model for Mobile Apps
type MobileAppDataSourceModel struct {
	AppId         types.String     `tfsdk:"app_id"`
	DisplayName   types.String     `tfsdk:"display_name"`
	Publisher     types.String     `tfsdk:"publisher"`
	Developer     types.String     `tfsdk:"developer"`
	Category      types.String     `tfsdk:"category"`
	ListAll       types.Bool       `tfsdk:"list_all"`
	ODataQuery    types.String     `tfsdk:"odata_query"`
	AppTypeFilter types.String     `tfsdk:"app_type_filter"`
	Items         []MobileAppModel `tfsdk:"items"`
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
	UploadState           types.Int32    `tfsdk:"upload_state"`
	DependentAppCount     types.Int32    `tfsdk:"dependent_app_count"`
	SupersededAppCount    types.Int32    `tfsdk:"superseded_app_count"`
	SupersedingAppCount   types.Int32    `tfsdk:"superseding_app_count"`
	RoleScopeTagIds       []types.String `tfsdk:"role_scope_tag_ids"`
	Categories            []types.String `tfsdk:"categories"`
}
