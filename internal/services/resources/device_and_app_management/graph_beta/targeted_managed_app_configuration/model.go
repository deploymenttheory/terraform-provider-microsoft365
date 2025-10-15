// Package graphBetaMobileAppConfiguration provides the resource implementation for Microsoft Graph Beta Mobile App Configuration
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-targetedmanagedappconfiguration?view=graph-rest-beta
package graphBetaTargetedManagedAppConfigurations

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TargetedManagedAppConfigurationResourceModel holds the configuration for a Mobile App Configuration.
type TargetedManagedAppConfigurationResourceModel struct {
	ID                          types.String                             `tfsdk:"id"`
	DisplayName                 types.String                             `tfsdk:"display_name"`
	Description                 types.String                             `tfsdk:"description"`
	CreatedDateTime             types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime        types.String                             `tfsdk:"last_modified_date_time"`
	Version                     types.String                             `tfsdk:"version"`
	RoleScopeTagIds             types.Set                                `tfsdk:"role_scope_tag_ids"`
	CustomSettings              []KeyValuePairResourceModel              `tfsdk:"custom_settings"`
	AppGroupType                types.String                             `tfsdk:"app_group_type"`
	Apps                        types.Set                                `tfsdk:"apps"`
	Assignments                 types.Set                                `tfsdk:"assignments"`
	DeployedAppCount            types.Int32                              `tfsdk:"deployed_app_count"`
	IsAssigned                  types.Bool                               `tfsdk:"is_assigned"`
	TargetedAppManagementLevels types.String                             `tfsdk:"targeted_app_management_levels"`
	SettingsCatalog             *DeviceConfigV2GraphServiceResourceModel `tfsdk:"settings_catalog"`
	Timeouts                    timeouts.Value                           `tfsdk:"timeouts"`
}

// KeyValuePairResourceModel represents a key-value pair configuration
type KeyValuePairResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// ManagedMobileAppResourceModel represents a managed mobile app
type ManagedMobileAppResourceModel struct {
	MobileAppIdentifier *MobileAppIdentifierModel `tfsdk:"mobile_app_identifier"`
	Version             types.String              `tfsdk:"version"`
}

// MobileAppIdentifierModel represents a mobile app identifier
type MobileAppIdentifierModel struct {
	Type         types.String `tfsdk:"type"`
	BundleId     types.String `tfsdk:"bundle_id"`
	PackageId    types.String `tfsdk:"package_id"`
	WindowsAppId types.String `tfsdk:"windows_app_id"`
}
