// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetapp?view=graph-rest-beta
package graphBetaWinGetApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WinGetAppResourceModel represents the Terraform resource model for a WinGetApp
type WinGetAppResourceModel struct {
	ID                    types.String                                   `tfsdk:"id"`
	DisplayName           types.String                                   `tfsdk:"display_name"`
	Description           types.String                                   `tfsdk:"description"`
	Publisher             types.String                                   `tfsdk:"publisher"`
	LargeIcon             types.Object                                   `tfsdk:"large_icon"`
	CreatedDateTime       types.String                                   `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String                                   `tfsdk:"last_modified_date_time"`
	IsFeatured            types.Bool                                     `tfsdk:"is_featured"`
	PrivacyInformationUrl types.String                                   `tfsdk:"privacy_information_url"`
	InformationUrl        types.String                                   `tfsdk:"information_url"`
	Owner                 types.String                                   `tfsdk:"owner"`
	Developer             types.String                                   `tfsdk:"developer"`
	Notes                 types.String                                   `tfsdk:"notes"`
	UploadState           types.Int64                                    `tfsdk:"upload_state"`
	PublishingState       types.String                                   `tfsdk:"publishing_state"`
	IsAssigned            types.Bool                                     `tfsdk:"is_assigned"`
	RoleScopeTagIds       []types.String                                 `tfsdk:"role_scope_tag_ids"`
	DependentAppCount     types.Int64                                    `tfsdk:"dependent_app_count"`
	SupersedingAppCount   types.Int64                                    `tfsdk:"superseding_app_count"`
	SupersededAppCount    types.Int64                                    `tfsdk:"superseded_app_count"`
	ManifestHash          types.String                                   `tfsdk:"manifest_hash"`
	PackageIdentifier     types.String                                   `tfsdk:"package_identifier"`
	InstallExperience     *WinGetAppInstallExperienceResourceModel       `tfsdk:"install_experience"`
	Assignments           *sharedmodels.MobileAppAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts              timeouts.Value                                 `tfsdk:"timeouts"`
}

// WinGetAppInstallExperienceModel represents the install experience structure
type WinGetAppInstallExperienceResourceModel struct {
	RunAsAccount types.String `tfsdk:"run_as_account"`
}
