// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-iosmanagedappprotection
package graphBetaDeviceAndAppManagementIosManagedAppProtection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IosManagedAppProtectionResourceModel represents the Terraform resource model
// for an iOS managed app protection policy.
type IosManagedAppProtectionResourceModel struct {
	// Computed-only
	ID                   types.String `tfsdk:"id"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	Version              types.String `tfsdk:"version"`
	IsAssigned           types.Bool   `tfsdk:"is_assigned"`
	DeployedAppCount     types.Int32  `tfsdk:"deployed_app_count"`

	// Required
	DisplayName types.String `tfsdk:"display_name"`

	// Optional
	Description                             types.String `tfsdk:"description"`
	PeriodOfflineBeforeAccessCheck          types.String `tfsdk:"period_offline_before_access_check"`
	PeriodOnlineBeforeAccessCheck           types.String `tfsdk:"period_online_before_access_check"`
	AllowedInboundDataTransferSources       types.String `tfsdk:"allowed_inbound_data_transfer_sources"`
	AllowedOutboundDataTransferDestinations types.String `tfsdk:"allowed_outbound_data_transfer_destinations"`
	OrganizationalCredentialsRequired       types.Bool   `tfsdk:"organizational_credentials_required"`
	AllowedOutboundClipboardSharingLevel    types.String `tfsdk:"allowed_outbound_clipboard_sharing_level"`
	DataBackupBlocked                       types.Bool   `tfsdk:"data_backup_blocked"`
	DeviceComplianceRequired                types.Bool   `tfsdk:"device_compliance_required"`
	ManagedBrowserToOpenLinksRequired       types.Bool   `tfsdk:"managed_browser_to_open_links_required"`
	SaveAsBlocked                           types.Bool   `tfsdk:"save_as_blocked"`
	PeriodOfflineBeforeWipeIsEnforced       types.String `tfsdk:"period_offline_before_wipe_is_enforced"`
	PinRequired                             types.Bool   `tfsdk:"pin_required"`
	MaximumPinRetries                       types.Int32  `tfsdk:"maximum_pin_retries"`
	SimplePinBlocked                        types.Bool   `tfsdk:"simple_pin_blocked"`
	MinimumPinLength                        types.Int32  `tfsdk:"minimum_pin_length"`
	PinCharacterSet                         types.String `tfsdk:"pin_character_set"`
	PeriodBeforePinReset                    types.String `tfsdk:"period_before_pin_reset"`
	AllowedDataStorageLocations             types.List   `tfsdk:"allowed_data_storage_locations"`
	ContactSyncBlocked                      types.Bool   `tfsdk:"contact_sync_blocked"`
	PrintBlocked                            types.Bool   `tfsdk:"print_blocked"`
	FingerprintBlocked                      types.Bool   `tfsdk:"fingerprint_blocked"`
	DisableAppPinIfDevicePinIsSet           types.Bool   `tfsdk:"disable_app_pin_if_device_pin_is_set"`
	MinimumRequiredOsVersion                types.String `tfsdk:"minimum_required_os_version"`
	MinimumWarningOsVersion                 types.String `tfsdk:"minimum_warning_os_version"`
	MinimumRequiredAppVersion               types.String `tfsdk:"minimum_required_app_version"`
	MinimumWarningAppVersion                types.String `tfsdk:"minimum_warning_app_version"`
	ManagedBrowser                          types.String `tfsdk:"managed_browser"`

	// Optional 
	FaceIdBlocked                                 types.Bool   `tfsdk:"face_id_blocked"`
	CustomBrowserProtocol                         types.String `tfsdk:"custom_browser_protocol"`
	ThirdPartyKeyboardsBlocked                    types.Bool   `tfsdk:"third_party_keyboards_blocked"`
	FilterOpenInToOnlyManagedApps                 types.Bool   `tfsdk:"filter_open_in_to_only_managed_apps"`
	DisableProtectionOfManagedOutboundOpenInData  types.Bool   `tfsdk:"disable_protection_of_managed_outbound_open_in_data"`
	AppDataEncryptionType                         types.String `tfsdk:"app_data_encryption_type"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}