// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowswebapp?view=graph-rest-beta

package graphBetaWindowsWebApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsWebAppResourceModel represents the root Terraform resource model for Windows Web Clip applications
type WindowsWebAppResourceModel struct {
	ID                    types.String                             `tfsdk:"id"`
	DisplayName           types.String                             `tfsdk:"display_name"`
	Description           types.String                             `tfsdk:"description"`
	Publisher             types.String                             `tfsdk:"publisher"`
	AppIcon               *sharedmodels.MobileAppIconResourceModel `tfsdk:"app_icon"`
	CreatedDateTime       types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured            types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl        types.String                             `tfsdk:"information_url"`
	Owner                 types.String                             `tfsdk:"owner"`
	Developer             types.String                             `tfsdk:"developer"`
	Notes                 types.String                             `tfsdk:"notes"`
	UploadState           types.Int32                              `tfsdk:"upload_state"`
	PublishingState       types.String                             `tfsdk:"publishing_state"`
	IsAssigned            types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds       types.Set                                `tfsdk:"role_scope_tag_ids"`
	DependentAppCount     types.Int32                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount   types.Int32                              `tfsdk:"superseding_app_count"`
	SupersededAppCount    types.Int32                              `tfsdk:"superseded_app_count"`
	AppUrl                types.String                             `tfsdk:"app_url"`
	Categories            types.Set                                `tfsdk:"categories"`
	Timeouts              timeouts.Value                           `tfsdk:"timeouts"`
}
