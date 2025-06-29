// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanageddevice-list?view=graph-rest-beta
package graphBetaManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ManagedDeviceDataSourceModel represents the Terraform data source model for Windows Managed Devices
type ManagedDeviceDataSourceModel struct {
	FilterType  types.String        `tfsdk:"filter_type"`
	FilterValue types.String        `tfsdk:"filter_value"`
	Items       []ManagedDeviceItem `tfsdk:"items"`
	Timeouts    timeouts.Value      `tfsdk:"timeouts"`
}

// ManagedDeviceItem represents an individual Windows Managed Device
type ManagedDeviceItem struct {
	ID                                        types.String                               `tfsdk:"id"`
	UserID                                    types.String                               `tfsdk:"user_id"`
	DeviceName                                types.String                               `tfsdk:"device_name"`
	OperatingSystem                           types.String                               `tfsdk:"operating_system"`
	OSVersion                                 types.String                               `tfsdk:"os_version"`
	ComplianceState                           types.String                               `tfsdk:"compliance_state"`
	ManagementState                           types.String                               `tfsdk:"management_state"`
	LastSyncDateTime                          types.String                               `tfsdk:"last_sync_date_time"`
	EnrolledDateTime                          types.String                               `tfsdk:"enrolled_date_time"`
	SerialNumber                              types.String                               `tfsdk:"serial_number"`
	Model                                     types.String                               `tfsdk:"model"`
	Manufacturer                              types.String                               `tfsdk:"manufacturer"`
	JailBroken                                types.String                               `tfsdk:"jail_broken"`
	EasActivated                              types.Bool                                 `tfsdk:"eas_activated"`
	EasDeviceId                               types.String                               `tfsdk:"eas_device_id"`
	AzureADDeviceId                           types.String                               `tfsdk:"azure_ad_device_id"`
	DeviceType                                types.String                               `tfsdk:"device_type"`
	OwnerType                                 types.String                               `tfsdk:"owner_type"`
	ManagedDeviceOwnerType                    types.String                               `tfsdk:"managed_device_owner_type"`
	UserPrincipalName                         types.String                               `tfsdk:"user_principal_name"`
	PhoneNumber                               types.String                               `tfsdk:"phone_number"`
	EmailAddress                              types.String                               `tfsdk:"email_address"`
	DeviceCategoryDisplayName                 types.String                               `tfsdk:"device_category_display_name"`
	IsSupervised                              types.Bool                                 `tfsdk:"is_supervised"`
	IsEncrypted                               types.Bool                                 `tfsdk:"is_encrypted"`
	WiFiMacAddress                            types.String                               `tfsdk:"wi_fi_mac_address"`
	EthernetMacAddress                        types.String                               `tfsdk:"ethernet_mac_address"`
	PhysicalMemoryInBytes                     types.Int64                                `tfsdk:"physical_memory_in_bytes"`
	ProcessorArchitecture                     types.String                               `tfsdk:"processor_architecture"`
	Notes                                     types.String                               `tfsdk:"notes"`
	HardwareInformation                       *HardwareInformation                       `tfsdk:"hardware_information"`
	DeviceActionResults                       []DeviceActionResult                       `tfsdk:"device_action_results"`
	UsersLoggedOn                             []LoggedOnUser                             `tfsdk:"users_logged_on"`
	ConfigurationManagerClientEnabledFeatures *ConfigurationManagerClientEnabledFeatures `tfsdk:"configuration_manager_client_enabled_features"`
	DeviceHealthAttestationState              *DeviceHealthAttestationState              `tfsdk:"device_health_attestation_state"`
	ConfigurationManagerClientHealthState     *ConfigurationManagerClientHealthState     `tfsdk:"configuration_manager_client_health_state"`
	ConfigurationManagerClientInformation     *ConfigurationManagerClientInformation     `tfsdk:"configuration_manager_client_information"`
	ChromeOSDeviceInfo                        []ChromeOSDeviceProperty                   `tfsdk:"chrome_os_device_info"`
	DeviceIdentityAttestationDetail           *DeviceIdentityAttestationDetail           `tfsdk:"device_identity_attestation_detail"`
}

type HardwareInformation struct {
	SerialNumber           types.String  `tfsdk:"serial_number"`
	Manufacturer           types.String  `tfsdk:"manufacturer"`
	Model                  types.String  `tfsdk:"model"`
	PhoneNumber            types.String  `tfsdk:"phone_number"`
	WifiMac                types.String  `tfsdk:"wifi_mac"`
	IsSupervised           types.Bool    `tfsdk:"is_supervised"`
	IsEncrypted            types.Bool    `tfsdk:"is_encrypted"`
	BatteryLevelPercentage types.Float64 `tfsdk:"battery_level_percentage"`
}

type DeviceActionResult struct {
	ActionName          types.String `tfsdk:"action_name"`
	ActionState         types.String `tfsdk:"action_state"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	LastUpdatedDateTime types.String `tfsdk:"last_updated_date_time"`
}

type LoggedOnUser struct {
	UserId            types.String `tfsdk:"user_id"`
	LastLogOnDateTime types.String `tfsdk:"last_log_on_date_time"`
}

type ConfigurationManagerClientEnabledFeatures struct {
	Inventory                types.Bool `tfsdk:"inventory"`
	ModernApps               types.Bool `tfsdk:"modern_apps"`
	ResourceAccess           types.Bool `tfsdk:"resource_access"`
	DeviceConfiguration      types.Bool `tfsdk:"device_configuration"`
	CompliancePolicy         types.Bool `tfsdk:"compliance_policy"`
	WindowsUpdateForBusiness types.Bool `tfsdk:"windows_update_for_business"`
	EndpointProtection       types.Bool `tfsdk:"endpoint_protection"`
	OfficeApps               types.Bool `tfsdk:"office_apps"`
}

type DeviceHealthAttestationState struct {
	LastUpdateDateTime types.String `tfsdk:"last_update_date_time"`
	BitLockerStatus    types.String `tfsdk:"bit_locker_status"`
	SecureBoot         types.String `tfsdk:"secure_boot"`
}

type ConfigurationManagerClientHealthState struct {
	State            types.String `tfsdk:"state"`
	ErrorCode        types.Int64  `tfsdk:"error_code"`
	LastSyncDateTime types.String `tfsdk:"last_sync_date_time"`
}

type ConfigurationManagerClientInformation struct {
	ClientIdentifier types.String `tfsdk:"client_identifier"`
	IsBlocked        types.Bool   `tfsdk:"is_blocked"`
	ClientVersion    types.String `tfsdk:"client_version"`
}

type ChromeOSDeviceProperty struct {
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
	Updatable types.Bool   `tfsdk:"updatable"`
}

type DeviceIdentityAttestationDetail struct {
	DeviceIdentityAttestationStatus types.String `tfsdk:"device_identity_attestation_status"`
}
