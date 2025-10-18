// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanageddevice-list?view=graph-rest-beta
package graphBetaManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ManagedDeviceDataSourceModel represents the Terraform data source model for Windows Managed Devices
type ManagedDeviceDataSourceModel struct {
	FilterType   types.String                       `tfsdk:"filter_type"`
	FilterValue  types.String                       `tfsdk:"filter_value"`
	ODataFilter  types.String                       `tfsdk:"odata_filter"`
	ODataTop     types.Int32                        `tfsdk:"odata_top"`
	ODataSkip    types.Int32                        `tfsdk:"odata_skip"`
	ODataSelect  types.String                       `tfsdk:"odata_select"`
	ODataOrderBy types.String                       `tfsdk:"odata_orderby"`
	ODataCount   types.Bool                         `tfsdk:"odata_count"`
	ODataSearch  types.String                       `tfsdk:"odata_search"`
	ODataExpand  types.String                       `tfsdk:"odata_expand"`
	Items        []ManagedDeviceDeviceDataItemModel `tfsdk:"items"`
	Timeouts     timeouts.Value                     `tfsdk:"timeouts"`
}

// ManagedDeviceDeviceDataItemModel represents an individual Windows Managed Device, aligned with the provided JSON structure.
type ManagedDeviceDeviceDataItemModel struct {
	ID                                          types.String                               `tfsdk:"id"`
	UserId                                      types.String                               `tfsdk:"user_id"`
	DeviceName                                  types.String                               `tfsdk:"device_name"`
	HardwareInformation                         *ManagedDeviceHardwareInformation          `tfsdk:"hardware_information"`
	OwnerType                                   types.String                               `tfsdk:"owner_type"`
	ManagedDeviceOwnerType                      types.String                               `tfsdk:"managed_device_owner_type"`
	DeviceActionResults                         []DeviceActionResult                       `tfsdk:"device_action_results"`
	ManagementState                             types.String                               `tfsdk:"management_state"`
	EnrolledDateTime                            types.String                               `tfsdk:"enrolled_date_time"`
	LastSyncDateTime                            types.String                               `tfsdk:"last_sync_date_time"`
	ChassisType                                 types.String                               `tfsdk:"chassis_type"`
	OperatingSystem                             types.String                               `tfsdk:"operating_system"`
	DeviceType                                  types.String                               `tfsdk:"device_type"`
	ComplianceState                             types.String                               `tfsdk:"compliance_state"`
	JailBroken                                  types.String                               `tfsdk:"jail_broken"`
	ManagementAgent                             types.String                               `tfsdk:"management_agent"`
	OSVersion                                   types.String                               `tfsdk:"os_version"`
	EasActivated                                types.Bool                                 `tfsdk:"eas_activated"`
	EasDeviceId                                 types.String                               `tfsdk:"eas_device_id"`
	EasActivationDateTime                       types.String                               `tfsdk:"eas_activation_date_time"`
	AadRegistered                               types.Bool                                 `tfsdk:"aad_registered"`
	AzureADRegistered                           types.Bool                                 `tfsdk:"azure_ad_registered"`
	DeviceEnrollmentType                        types.String                               `tfsdk:"device_enrollment_type"`
	LostModeState                               types.String                               `tfsdk:"lost_mode_state"`
	ActivationLockBypassCode                    types.String                               `tfsdk:"activation_lock_bypass_code"`
	EmailAddress                                types.String                               `tfsdk:"email_address"`
	AzureActiveDirectoryDeviceId                types.String                               `tfsdk:"azure_active_directory_device_id"`
	AzureADDeviceId                             types.String                               `tfsdk:"azure_ad_device_id"`
	DeviceRegistrationState                     types.String                               `tfsdk:"device_registration_state"`
	DeviceCategoryDisplayName                   types.String                               `tfsdk:"device_category_display_name"`
	IsSupervised                                types.Bool                                 `tfsdk:"is_supervised"`
	ExchangeLastSuccessfulSyncDateTime          types.String                               `tfsdk:"exchange_last_successful_sync_date_time"`
	ExchangeAccessState                         types.String                               `tfsdk:"exchange_access_state"`
	ExchangeAccessStateReason                   types.String                               `tfsdk:"exchange_access_state_reason"`
	RemoteAssistanceSessionUrl                  types.String                               `tfsdk:"remote_assistance_session_url"`
	RemoteAssistanceSessionErrorDetails         types.String                               `tfsdk:"remote_assistance_session_error_details"`
	IsEncrypted                                 types.Bool                                 `tfsdk:"is_encrypted"`
	UserPrincipalName                           types.String                               `tfsdk:"user_principal_name"`
	Model                                       types.String                               `tfsdk:"model"`
	Manufacturer                                types.String                               `tfsdk:"manufacturer"`
	IMEI                                        types.String                               `tfsdk:"imei"`
	ComplianceGracePeriodExpirationDateTime     types.String                               `tfsdk:"compliance_grace_period_expiration_date_time"`
	SerialNumber                                types.String                               `tfsdk:"serial_number"`
	PhoneNumber                                 types.String                               `tfsdk:"phone_number"`
	AndroidSecurityPatchLevel                   types.String                               `tfsdk:"android_security_patch_level"`
	UserDisplayName                             types.String                               `tfsdk:"user_display_name"`
	ConfigurationManagerClientEnabledFeatures   *ConfigurationManagerClientEnabledFeatures `tfsdk:"configuration_manager_client_enabled_features"`
	WiFiMacAddress                              types.String                               `tfsdk:"wi_fi_mac_address"`
	DeviceHealthAttestationState                *DeviceHealthAttestationState              `tfsdk:"device_health_attestation_state"`
	SubscriberCarrier                           types.String                               `tfsdk:"subscriber_carrier"`
	MEID                                        types.String                               `tfsdk:"meid"`
	TotalStorageSpaceInBytes                    types.Int64                                `tfsdk:"total_storage_space_in_bytes"`
	FreeStorageSpaceInBytes                     types.Int64                                `tfsdk:"free_storage_space_in_bytes"`
	ManagedDeviceName                           types.String                               `tfsdk:"managed_device_name"`
	PartnerReportedThreatState                  types.String                               `tfsdk:"partner_reported_threat_state"`
	RetireAfterDateTime                         types.String                               `tfsdk:"retire_after_date_time"`
	UsersLoggedOn                               []LoggedOnUser                             `tfsdk:"users_logged_on"`
	PreferMdmOverGroupPolicyAppliedDateTime     types.String                               `tfsdk:"prefer_mdm_over_group_policy_applied_date_time"`
	AutopilotEnrolled                           types.Bool                                 `tfsdk:"autopilot_enrolled"`
	RequireUserEnrollmentApproval               types.Bool                                 `tfsdk:"require_user_enrollment_approval"`
	ManagementCertificateExpirationDate         types.String                               `tfsdk:"management_certificate_expiration_date"`
	ICCID                                       types.String                               `tfsdk:"iccid"`
	UDID                                        types.String                               `tfsdk:"udid"`
	RoleScopeTagIds                             []types.String                             `tfsdk:"role_scope_tag_ids"`
	WindowsActiveMalwareCount                   types.Int64                                `tfsdk:"windows_active_malware_count"`
	WindowsRemediatedMalwareCount               types.Int64                                `tfsdk:"windows_remediated_malware_count"`
	Notes                                       types.String                               `tfsdk:"notes"`
	ConfigurationManagerClientHealthState       *ConfigurationManagerClientHealthState     `tfsdk:"configuration_manager_client_health_state"`
	ConfigurationManagerClientInformation       *ConfigurationManagerClientInformation     `tfsdk:"configuration_manager_client_information"`
	EthernetMacAddress                          types.String                               `tfsdk:"ethernet_mac_address"`
	PhysicalMemoryInBytes                       types.Int64                                `tfsdk:"physical_memory_in_bytes"`
	ProcessorArchitecture                       types.String                               `tfsdk:"processor_architecture"`
	SpecificationVersion                        types.String                               `tfsdk:"specification_version"`
	JoinType                                    types.String                               `tfsdk:"join_type"`
	SkuFamily                                   types.String                               `tfsdk:"sku_family"`
	SecurityPatchLevel                          types.String                               `tfsdk:"security_patch_level"`
	SkuNumber                                   types.Int32                                `tfsdk:"sku_number"`
	ManagementFeatures                          types.String                               `tfsdk:"management_features"`
	ChromeOSDeviceInfo                          []ChromeOSDeviceInfo                       `tfsdk:"chrome_os_device_info"`
	EnrollmentProfileName                       types.String                               `tfsdk:"enrollment_profile_name"`
	BootstrapTokenEscrowed                      types.Bool                                 `tfsdk:"bootstrap_token_escrowed"`
	DeviceFirmwareConfigurationInterfaceManaged types.Bool                                 `tfsdk:"device_firmware_configuration_interface_managed"`
	DeviceIdentityAttestationDetail             *DeviceIdentityAttestationDetail           `tfsdk:"device_identity_attestation_detail"`
}

// DeviceActionResult represents an item in deviceActionResults
type DeviceActionResult struct {
	ActionName          types.String `tfsdk:"action_name"`
	ActionState         types.String `tfsdk:"action_state"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	LastUpdatedDateTime types.String `tfsdk:"last_updated_date_time"`
}

// LoggedOnUser represents an item in usersLoggedOn
type LoggedOnUser struct {
	UserId            types.String `tfsdk:"user_id"`
	LastLogOnDateTime types.String `tfsdk:"last_log_on_date_time"`
}

// ChromeOSDeviceInfo represents an item in chromeOSDeviceInfo
type ChromeOSDeviceInfo struct {
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
	Updatable types.Bool   `tfsdk:"updatable"`
}

// DeviceHealthAttestationState represents the deviceHealthAttestationState object
type DeviceHealthAttestationState struct {
	LastUpdateDateTime                       types.String `tfsdk:"last_update_date_time"`
	ContentNamespaceUrl                      types.String `tfsdk:"content_namespace_url"`
	DeviceHealthAttestationStatus            types.String `tfsdk:"device_health_attestation_status"`
	ContentVersion                           types.String `tfsdk:"content_version"`
	IssuedDateTime                           types.String `tfsdk:"issued_date_time"`
	AttestationIdentityKey                   types.String `tfsdk:"attestation_identity_key"`
	ResetCount                               types.Int64  `tfsdk:"reset_count"`
	RestartCount                             types.Int64  `tfsdk:"restart_count"`
	DataExcutionPolicy                       types.String `tfsdk:"data_excution_policy"`
	BitLockerStatus                          types.String `tfsdk:"bit_locker_status"`
	BootManagerVersion                       types.String `tfsdk:"boot_manager_version"`
	CodeIntegrityCheckVersion                types.String `tfsdk:"code_integrity_check_version"`
	SecureBoot                               types.String `tfsdk:"secure_boot"`
	BootDebugging                            types.String `tfsdk:"boot_debugging"`
	OperatingSystemKernelDebugging           types.String `tfsdk:"operating_system_kernel_debugging"`
	CodeIntegrity                            types.String `tfsdk:"code_integrity"`
	TestSigning                              types.String `tfsdk:"test_signing"`
	SafeMode                                 types.String `tfsdk:"safe_mode"`
	WindowsPE                                types.String `tfsdk:"windows_pe"`
	EarlyLaunchAntiMalwareDriverProtection   types.String `tfsdk:"early_launch_anti_malware_driver_protection"`
	VirtualSecureMode                        types.String `tfsdk:"virtual_secure_mode"`
	PcrHashAlgorithm                         types.String `tfsdk:"pcr_hash_algorithm"`
	BootAppSecurityVersion                   types.String `tfsdk:"boot_app_security_version"`
	BootManagerSecurityVersion               types.String `tfsdk:"boot_manager_security_version"`
	TpmVersion                               types.String `tfsdk:"tpm_version"`
	Pcr0                                     types.String `tfsdk:"pcr0"`
	SecureBootConfigurationPolicyFingerPrint types.String `tfsdk:"secure_boot_configuration_policy_finger_print"`
	CodeIntegrityPolicy                      types.String `tfsdk:"code_integrity_policy"`
	BootRevisionListInfo                     types.String `tfsdk:"boot_revision_list_info"`
	OperatingSystemRevListInfo               types.String `tfsdk:"operating_system_rev_list_info"`
	HealthStatusMismatchInfo                 types.String `tfsdk:"health_status_mismatch_info"`
	HealthAttestationSupportedStatus         types.String `tfsdk:"health_attestation_supported_status"`
	MemoryIntegrityProtection                types.String `tfsdk:"memory_integrity_protection"`
	MemoryAccessProtection                   types.String `tfsdk:"memory_access_protection"`
	VirtualizationBasedSecurity              types.String `tfsdk:"virtualization_based_security"`
	FirmwareProtection                       types.String `tfsdk:"firmware_protection"`
	SystemManagementMode                     types.String `tfsdk:"system_management_mode"`
	SecuredCorePC                            types.String `tfsdk:"secured_core_pc"`
}

// ConfigurationManagerClientEnabledFeatures represents the configurationManagerClientEnabledFeatures object
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

// ConfigurationManagerClientHealthState represents the configurationManagerClientHealthState object
type ConfigurationManagerClientHealthState struct {
	State            types.String `tfsdk:"state"`
	ErrorCode        types.Int64  `tfsdk:"error_code"`
	LastSyncDateTime types.String `tfsdk:"last_sync_date_time"`
}

// ConfigurationManagerClientInformation represents the configurationManagerClientInformation object
type ConfigurationManagerClientInformation struct {
	ClientIdentifier types.String `tfsdk:"client_identifier"`
	IsBlocked        types.Bool   `tfsdk:"is_blocked"`
	ClientVersion    types.String `tfsdk:"client_version"`
}

// DeviceIdentityAttestationDetail represents the deviceIdentityAttestationDetail object
type DeviceIdentityAttestationDetail struct {
	DeviceIdentityAttestationStatus types.String `tfsdk:"device_identity_attestation_status"`
}

// ManagedDeviceHardwareInformation represents the hardwareInformation nested object in the JSON.
type ManagedDeviceHardwareInformation struct {
	SerialNumber                                                   types.String            `tfsdk:"serial_number"`
	TotalStorageSpace                                              types.Int64             `tfsdk:"total_storage_space"`
	FreeStorageSpace                                               types.Int64             `tfsdk:"free_storage_space"`
	IMEI                                                           types.String            `tfsdk:"imei"`
	MEID                                                           types.String            `tfsdk:"meid"`
	Manufacturer                                                   types.String            `tfsdk:"manufacturer"`
	Model                                                          types.String            `tfsdk:"model"`
	PhoneNumber                                                    types.String            `tfsdk:"phone_number"`
	SubscriberCarrier                                              types.String            `tfsdk:"subscriber_carrier"`
	CellularTechnology                                             types.String            `tfsdk:"cellular_technology"`
	WifiMac                                                        types.String            `tfsdk:"wifi_mac"`
	OperatingSystemLanguage                                        types.String            `tfsdk:"operating_system_language"`
	IsSupervised                                                   types.Bool              `tfsdk:"is_supervised"`
	IsEncrypted                                                    types.Bool              `tfsdk:"is_encrypted"`
	BatterySerialNumber                                            types.String            `tfsdk:"battery_serial_number"`
	BatteryHealthPercentage                                        types.Int64             `tfsdk:"battery_health_percentage"`
	BatteryChargeCycles                                            types.Int64             `tfsdk:"battery_charge_cycles"`
	IsSharedDevice                                                 types.Bool              `tfsdk:"is_shared_device"`
	SharedDeviceCachedUsers                                        []SharedAppleDeviceUser `tfsdk:"shared_device_cached_users"`
	TPMSpecificationVersion                                        types.String            `tfsdk:"tpm_specification_version"`
	OperatingSystemEdition                                         types.String            `tfsdk:"operating_system_edition"`
	DeviceFullQualifiedDomainName                                  types.String            `tfsdk:"device_full_qualified_domain_name"`
	DeviceGuardVirtualizationBasedSecurityHardwareRequirementState types.String            `tfsdk:"device_guard_virtualization_based_security_hardware_requirement_state"`
	DeviceGuardVirtualizationBasedSecurityState                    types.String            `tfsdk:"device_guard_virtualization_based_security_state"`
	DeviceGuardLocalSystemAuthorityCredentialGuardState            types.String            `tfsdk:"device_guard_local_system_authority_credential_guard_state"`
	OSBuildNumber                                                  types.String            `tfsdk:"os_build_number"`
	OperatingSystemProductType                                     types.Int64             `tfsdk:"operating_system_product_type"`
	IPAddressV4                                                    types.String            `tfsdk:"ip_address_v4"`
	SubnetAddress                                                  types.String            `tfsdk:"subnet_address"`
	ESIMIdentifier                                                 types.String            `tfsdk:"esim_identifier"`
	SystemManagementBIOSVersion                                    types.String            `tfsdk:"system_management_bios_version"`
	TPMManufacturer                                                types.String            `tfsdk:"tpm_manufacturer"`
	TPMVersion                                                     types.String            `tfsdk:"tpm_version"`
	WiredIPv4Addresses                                             []types.String          `tfsdk:"wired_ipv4_addresses"`
	BatteryLevelPercentage                                         types.Float64           `tfsdk:"battery_level_percentage"`
	ResidentUsersCount                                             types.Int64             `tfsdk:"resident_users_count"`
	ProductName                                                    types.String            `tfsdk:"product_name"`
	DeviceLicensingStatus                                          types.String            `tfsdk:"device_licensing_status"`
	DeviceLicensingLastErrorCode                                   types.Int64             `tfsdk:"device_licensing_last_error_code"`
	DeviceLicensingLastErrorDescription                            types.String            `tfsdk:"device_licensing_last_error_description"`
}

// SharedAppleDeviceUser represents an item in sharedDeviceCachedUsers
type SharedAppleDeviceUser struct {
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
	DataToSync        types.Bool   `tfsdk:"data_to_sync"`
	DataQuota         types.Int64  `tfsdk:"data_quota"`
	DataUsed          types.Int64  `tfsdk:"data_used"`
}
