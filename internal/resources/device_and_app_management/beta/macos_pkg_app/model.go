// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta
package graphBetaMacosPkgApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MacOSPkgAppResourceModel struct {
	ID                              types.String                             `tfsdk:"id"`
	DisplayName                     types.String                             `tfsdk:"display_name"`
	Description                     types.String                             `tfsdk:"description"`
	Publisher                       types.String                             `tfsdk:"publisher"`
	LargeIcon                       sharedmodels.MimeContentResourceModel    `tfsdk:"large_icon"`
	CreatedDateTime                 types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime            types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured                      types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl           types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl                  types.String                             `tfsdk:"information_url"`
	Owner                           types.String                             `tfsdk:"owner"`
	Developer                       types.String                             `tfsdk:"developer"`
	Notes                           types.String                             `tfsdk:"notes"`
	UploadState                     types.Int64                              `tfsdk:"upload_state"`
	PublishingState                 types.String                             `tfsdk:"publishing_state"`
	IsAssigned                      types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds                 []types.String                           `tfsdk:"role_scope_tag_ids"`
	DependentAppCount               types.Int64                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount             types.Int64                              `tfsdk:"superseding_app_count"`
	SupersededAppCount              types.Int64                              `tfsdk:"superseded_app_count"`
	CommittedContentVersion         types.String                             `tfsdk:"committed_content_version"`
	FileName                        types.String                             `tfsdk:"file_name"`
	Size                            types.Int64                              `tfsdk:"size"`
	PrimaryBundleId                 types.String                             `tfsdk:"primary_bundle_id"`
	PrimaryBundleVersion            types.String                             `tfsdk:"primary_bundle_version"`
	IncludedApps                    []MacOSIncludedAppResourceModel          `tfsdk:"included_apps"`
	IgnoreVersionDetection          types.Bool                               `tfsdk:"ignore_version_detection"`
	MinimumSupportedOperatingSystem MacOSMinimumOperatingSystemResourceModel `tfsdk:"minimum_supported_operating_system"`
	PreInstallScript                MacOSAppScriptResourceModel              `tfsdk:"pre_install_script"`
	PostInstallScript               MacOSAppScriptResourceModel              `tfsdk:"post_install_script"`
	Timeouts                        timeouts.Value                           `tfsdk:"timeouts"`
}

type MacOSIncludedAppResourceModel struct {
	BundleId      types.String `tfsdk:"bundle_id"`
	BundleVersion types.String `tfsdk:"bundle_version"`
}

type MacOSMinimumOperatingSystemResourceModel struct {
	V10_7  types.Bool `tfsdk:"v10_7"`
	V10_8  types.Bool `tfsdk:"v10_8"`
	V10_9  types.Bool `tfsdk:"v10_9"`
	V10_10 types.Bool `tfsdk:"v10_10"`
	V10_11 types.Bool `tfsdk:"v10_11"`
	V10_12 types.Bool `tfsdk:"v10_12"`
	V10_13 types.Bool `tfsdk:"v10_13"`
	V10_14 types.Bool `tfsdk:"v10_14"`
	V10_15 types.Bool `tfsdk:"v10_15"`
	V11_0  types.Bool `tfsdk:"v11_0"`
	V12_0  types.Bool `tfsdk:"v12_0"`
	V13_0  types.Bool `tfsdk:"v13_0"`
	V14_0  types.Bool `tfsdk:"v14_0"`
}

type MacOSAppScriptResourceModel struct {
	ScriptContent types.String `tfsdk:"script_content"`
}
