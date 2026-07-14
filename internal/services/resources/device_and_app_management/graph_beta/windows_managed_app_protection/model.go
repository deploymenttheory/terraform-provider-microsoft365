// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-windowsmanagedappprotection
package graphBetaDeviceAndAppManagementWindowsManagedAppProtection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsManagedAppProtectionResourceModel represents the Terraform resource model
// for a Windows managed app protection policy.
type WindowsManagedAppProtectionResourceModel struct {
	// Computed-only — set by the API, never writable
	ID                   types.String `tfsdk:"id"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	Version              types.String `tfsdk:"version"`
	IsAssigned           types.Bool   `tfsdk:"is_assigned"`
	DeployedAppCount     types.Int64  `tfsdk:"deployed_app_count"`

	// Required
	DisplayName types.String `tfsdk:"display_name"`

	// Optional
	Description                              types.String `tfsdk:"description"`
	RoleScopeTagIds                          types.List   `tfsdk:"role_scope_tag_ids"`
	PrintBlocked                             types.Bool   `tfsdk:"print_blocked"`
	AllowedInboundDataTransferSources        types.String `tfsdk:"allowed_inbound_data_transfer_sources"`
	AllowedOutboundClipboardSharingLevel     types.String `tfsdk:"allowed_outbound_clipboard_sharing_level"`
	AllowedOutboundDataTransferDestinations  types.String `tfsdk:"allowed_outbound_data_transfer_destinations"`
	AppActionIfUnableToAuthenticateUser      types.String `tfsdk:"app_action_if_unable_to_authenticate_user"`
	MaximumAllowedDeviceThreatLevel          types.String `tfsdk:"maximum_allowed_device_threat_level"`
	MobileThreatDefenseRemediationAction     types.String `tfsdk:"mobile_threat_defense_remediation_action"`
	MinimumRequiredSdkVersion                types.String `tfsdk:"minimum_required_sdk_version"`
	MinimumWipeSdkVersion                    types.String `tfsdk:"minimum_wipe_sdk_version"`
	MinimumRequiredOsVersion                 types.String `tfsdk:"minimum_required_os_version"`
	MinimumWarningOsVersion                  types.String `tfsdk:"minimum_warning_os_version"`
	MinimumWipeOsVersion                     types.String `tfsdk:"minimum_wipe_os_version"`
	MinimumRequiredAppVersion                types.String `tfsdk:"minimum_required_app_version"`
	MinimumWarningAppVersion                 types.String `tfsdk:"minimum_warning_app_version"`
	MinimumWipeAppVersion                    types.String `tfsdk:"minimum_wipe_app_version"`
	MaximumRequiredOsVersion                 types.String `tfsdk:"maximum_required_os_version"`
	MaximumWarningOsVersion                  types.String `tfsdk:"maximum_warning_os_version"`
	MaximumWipeOsVersion                     types.String `tfsdk:"maximum_wipe_os_version"`
	PeriodOfflineBeforeWipeIsEnforced        types.String `tfsdk:"period_offline_before_wipe_is_enforced"`
	PeriodOfflineBeforeAccessCheck           types.String `tfsdk:"period_offline_before_access_check"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
