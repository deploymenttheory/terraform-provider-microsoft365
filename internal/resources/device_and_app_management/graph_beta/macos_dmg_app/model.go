// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosdmgapp?view=graph-rest-beta

package graphBetaMacOSDmgApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSDmgAppResourceModel represents the root Terraform resource model for intune macOS DMG applications
type MacOSDmgAppResourceModel struct {
	ID                    types.String                             `tfsdk:"id"`
	DisplayName           types.String                             `tfsdk:"display_name"`
	Description           types.String                             `tfsdk:"description"`
	Publisher             types.String                             `tfsdk:"publisher"`
	AppIcon               *sharedmodels.MobileAppIconResourceModel `tfsdk:"app_icon"`
	CreatedDateTime       types.String                             `tfsdk:"created_date_time"`
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
	Categories            types.Set                                `tfsdk:"categories"`
	Relationships         types.List                               `tfsdk:"relationships"`
	MacOSDmgApp           *MacOSDmgAppDetailsResourceModel         `tfsdk:"macos_dmg_app"`
	AppInstaller          types.Object                             `tfsdk:"app_installer"`
	ContentVersion        types.List                               `tfsdk:"content_version"`
	Timeouts              timeouts.Value                           `tfsdk:"timeouts"`
}

// MobileAppRelationshipResourceModel represents the Terraform resource model for a Mobile App Relationship
type MobileAppRelationshipResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	SourceDisplayName          types.String `tfsdk:"source_display_name"`
	SourceDisplayVersion       types.String `tfsdk:"source_display_version"`
	SourceId                   types.String `tfsdk:"source_id"`
	SourcePublisherDisplayName types.String `tfsdk:"source_publisher_display_name"`
	TargetDisplayName          types.String `tfsdk:"target_display_name"`
	TargetDisplayVersion       types.String `tfsdk:"target_display_version"`
	TargetId                   types.String `tfsdk:"target_id"`
	TargetPublisher            types.String `tfsdk:"target_publisher"`
	TargetPublisherDisplayName types.String `tfsdk:"target_publisher_display_name"`
	TargetType                 types.String `tfsdk:"target_type"`
}

// MacOSDmgAppDetailsResourceModel represents the Terraform resource model for a MacOS DMG Application
// Based on Microsoft Graph MicrosoftGraphMacOSDmgApp which extends MicrosoftGraphMobileLobApp
type MacOSDmgAppDetailsResourceModel struct {
	IgnoreVersionDetection          types.Bool                                `tfsdk:"ignore_version_detection"`
	IncludedApps                    types.Set                                 `tfsdk:"included_apps"`
	MinimumSupportedOperatingSystem *MacOSMinimumOperatingSystemResourceModel `tfsdk:"minimum_supported_operating_system"`
	PrimaryBundleId                 types.String                              `tfsdk:"primary_bundle_id"`
	PrimaryBundleVersion            types.String                              `tfsdk:"primary_bundle_version"`
}

// MacOSMinimumOperatingSystemResourceModel represents the minimum OS requirements for macOS
type MacOSMinimumOperatingSystemResourceModel struct {
	V107  types.Bool `tfsdk:"v10_7"`  // OS X 10.7 or later
	V108  types.Bool `tfsdk:"v10_8"`  // OS X 10.8 or later
	V109  types.Bool `tfsdk:"v10_9"`  // OS X 10.9 or later
	V1010 types.Bool `tfsdk:"v10_10"` // OS X 10.10 or later
	V1011 types.Bool `tfsdk:"v10_11"` // OS X 10.11 or later
	V1012 types.Bool `tfsdk:"v10_12"` // macOS 10.12 or later
	V1013 types.Bool `tfsdk:"v10_13"` // macOS 10.13 or later
	V1014 types.Bool `tfsdk:"v10_14"` // macOS 10.14 or later
	V1015 types.Bool `tfsdk:"v10_15"` // macOS 10.15 or later
	V110  types.Bool `tfsdk:"v11_0"`  // macOS 11.0 or later
	V120  types.Bool `tfsdk:"v12_0"`  // macOS 12.0 or later
	V130  types.Bool `tfsdk:"v13_0"`  // macOS 13.0 or later
	V140  types.Bool `tfsdk:"v14_0"`  // macOS 14.0 or later
	V150  types.Bool `tfsdk:"v15_0"`  // macOS 15.0 or later
}

// MacOSIncludedAppResourceModel represents an included app in the DMG
// This uses the same structure as PKG apps for IncludedApps
type MacOSIncludedAppResourceModel struct {
	BundleId      types.String `tfsdk:"bundle_id"`
	BundleVersion types.String `tfsdk:"bundle_version"`
}
