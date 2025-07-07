// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managediosstoreapp?view=graph-rest-beta

package graphBetaIOSStoreApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IOSStoreAppResourceModel represents the root Terraform resource model for iOS Store applications
type IOSStoreAppResourceModel struct {
	ID                              types.String                             `tfsdk:"id"`
	DisplayName                     types.String                             `tfsdk:"display_name"`
	Description                     types.String                             `tfsdk:"description"`
	Publisher                       types.String                             `tfsdk:"publisher"`
	AppIcon                         *sharedmodels.MobileAppIconResourceModel `tfsdk:"app_icon"`
	CreatedDateTime                 types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime            types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured                      types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl           types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl                  types.String                             `tfsdk:"information_url"`
	Owner                           types.String                             `tfsdk:"owner"`
	Developer                       types.String                             `tfsdk:"developer"`
	Notes                           types.String                             `tfsdk:"notes"`
	UploadState                     types.Int32                              `tfsdk:"upload_state"`
	PublishingState                 types.String                             `tfsdk:"publishing_state"`
	IsAssigned                      types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds                 types.Set                                `tfsdk:"role_scope_tag_ids"`
	DependentAppCount               types.Int32                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount             types.Int32                              `tfsdk:"superseding_app_count"`
	SupersededAppCount              types.Int32                              `tfsdk:"superseded_app_count"`
	AppStoreUrl                     types.String                             `tfsdk:"app_store_url"`
	ApplicableDeviceType            *IOSDeviceTypeResourceModel              `tfsdk:"applicable_device_type"`
	MinimumSupportedOperatingSystem *IOSMinimumOperatingSystemResourceModel  `tfsdk:"minimum_supported_operating_system"`
	Categories                      types.Set                                `tfsdk:"categories"`
	Relationships                   types.List                               `tfsdk:"relationships"`
	Timeouts                        timeouts.Value                           `tfsdk:"timeouts"`
}

// IOSDeviceTypeResourceModel represents the iOS device type compatibility
type IOSDeviceTypeResourceModel struct {
	IPad          types.Bool `tfsdk:"ipad"`
	IPhoneAndIPod types.Bool `tfsdk:"iphone_and_ipod"`
}

// IOSMinimumOperatingSystemResourceModel represents the minimum supported iOS version
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta
type IOSMinimumOperatingSystemResourceModel struct {
	V8_0  types.Bool `tfsdk:"v8_0"`
	V9_0  types.Bool `tfsdk:"v9_0"`
	V10_0 types.Bool `tfsdk:"v10_0"`
	V11_0 types.Bool `tfsdk:"v11_0"`
	V12_0 types.Bool `tfsdk:"v12_0"`
	V13_0 types.Bool `tfsdk:"v13_0"`
	V14_0 types.Bool `tfsdk:"v14_0"`
	V15_0 types.Bool `tfsdk:"v15_0"`
	V16_0 types.Bool `tfsdk:"v16_0"`
	V17_0 types.Bool `tfsdk:"v17_0"`
	V18_0 types.Bool `tfsdk:"v18_0"`
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

// IOSStoreAppAssignmentResourceModel represents the Terraform resource model for app assignments
type IOSStoreAppAssignmentResourceModel struct {
	ID       types.String                        `tfsdk:"id"`
	Intent   types.String                        `tfsdk:"intent"`
	Target   types.Object                        `tfsdk:"target"`
	Settings *IOSStoreAppAssignmentSettingsModel `tfsdk:"settings"`
	Source   types.String                        `tfsdk:"source"`
	SourceId types.String                        `tfsdk:"source_id"`
}

// IOSStoreAppAssignmentSettingsModel represents the iOS Store App assignment settings
type IOSStoreAppAssignmentSettingsModel struct {
	VpnConfigurationId       types.String `tfsdk:"vpn_configuration_id"`
	UninstallOnDeviceRemoval types.Bool   `tfsdk:"uninstall_on_device_removal"`
	IsRemovable              types.Bool   `tfsdk:"is_removable"`
	PreventManagedAppBackup  types.Bool   `tfsdk:"prevent_managed_app_backup"`
}
