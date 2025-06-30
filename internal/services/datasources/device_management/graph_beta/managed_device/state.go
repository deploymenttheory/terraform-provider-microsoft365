package graphBetaManagedDevice

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Managed Device Graph API object to the Terraform model
func MapRemoteStateToDataSource(data graphmodels.ManagedDeviceable) ManagedDeviceDeviceDataItemModel {
	return ManagedDeviceDeviceDataItemModel{
		ID:                                      convert.GraphToFrameworkString(data.GetId()),
		UserId:                                  convert.GraphToFrameworkString(data.GetUserId()),
		DeviceName:                              convert.GraphToFrameworkString(data.GetDeviceName()),
		OwnerType:                               convert.GraphToFrameworkEnum(data.GetOwnerType()),
		ManagedDeviceOwnerType:                  convert.GraphToFrameworkEnum(data.GetManagedDeviceOwnerType()),
		ManagementState:                         convert.GraphToFrameworkEnum(data.GetManagementState()),
		EnrolledDateTime:                        convert.GraphToFrameworkTime(data.GetEnrolledDateTime()),
		LastSyncDateTime:                        convert.GraphToFrameworkTime(data.GetLastSyncDateTime()),
		ChassisType:                             convert.GraphToFrameworkEnum(data.GetChassisType()),
		OperatingSystem:                         convert.GraphToFrameworkString(data.GetOperatingSystem()),
		DeviceType:                              convert.GraphToFrameworkEnum(data.GetDeviceType()),
		ComplianceState:                         convert.GraphToFrameworkEnum(data.GetComplianceState()),
		JailBroken:                              convert.GraphToFrameworkString(data.GetJailBroken()),
		ManagementAgent:                         convert.GraphToFrameworkEnum(data.GetManagementAgent()),
		OSVersion:                               convert.GraphToFrameworkString(data.GetOsVersion()),
		EasActivated:                            convert.GraphToFrameworkBool(data.GetEasActivated()),
		EasDeviceId:                             convert.GraphToFrameworkString(data.GetEasDeviceId()),
		EasActivationDateTime:                   convert.GraphToFrameworkTime(data.GetEasActivationDateTime()),
		AadRegistered:                           convert.GraphToFrameworkBool(data.GetAadRegistered()),
		AzureADRegistered:                       convert.GraphToFrameworkBool(data.GetAzureADRegistered()),
		DeviceEnrollmentType:                    convert.GraphToFrameworkEnum(data.GetDeviceEnrollmentType()),
		LostModeState:                           convert.GraphToFrameworkEnum(data.GetLostModeState()),
		ActivationLockBypassCode:                convert.GraphToFrameworkString(data.GetActivationLockBypassCode()),
		EmailAddress:                            convert.GraphToFrameworkString(data.GetEmailAddress()),
		AzureActiveDirectoryDeviceId:            convert.GraphToFrameworkString(data.GetAzureActiveDirectoryDeviceId()),
		AzureADDeviceId:                         convert.GraphToFrameworkString(data.GetAzureADDeviceId()),
		DeviceRegistrationState:                 convert.GraphToFrameworkEnum(data.GetDeviceRegistrationState()),
		DeviceCategoryDisplayName:               convert.GraphToFrameworkString(data.GetDeviceCategoryDisplayName()),
		IsSupervised:                            convert.GraphToFrameworkBool(data.GetIsSupervised()),
		ExchangeLastSuccessfulSyncDateTime:      convert.GraphToFrameworkTime(data.GetExchangeLastSuccessfulSyncDateTime()),
		ExchangeAccessState:                     convert.GraphToFrameworkEnum(data.GetExchangeAccessState()),
		ExchangeAccessStateReason:               convert.GraphToFrameworkEnum(data.GetExchangeAccessStateReason()),
		RemoteAssistanceSessionUrl:              convert.GraphToFrameworkString(data.GetRemoteAssistanceSessionUrl()),
		RemoteAssistanceSessionErrorDetails:     convert.GraphToFrameworkString(data.GetRemoteAssistanceSessionErrorDetails()),
		IsEncrypted:                             convert.GraphToFrameworkBool(data.GetIsEncrypted()),
		UserPrincipalName:                       convert.GraphToFrameworkString(data.GetUserPrincipalName()),
		Model:                                   convert.GraphToFrameworkString(data.GetModel()),
		Manufacturer:                            convert.GraphToFrameworkString(data.GetManufacturer()),
		IMEI:                                    convert.GraphToFrameworkString(data.GetImei()),
		ComplianceGracePeriodExpirationDateTime: convert.GraphToFrameworkTime(data.GetComplianceGracePeriodExpirationDateTime()),
		SerialNumber:                            convert.GraphToFrameworkString(data.GetSerialNumber()),
		PhoneNumber:                             convert.GraphToFrameworkString(data.GetPhoneNumber()),
		AndroidSecurityPatchLevel:               convert.GraphToFrameworkString(data.GetAndroidSecurityPatchLevel()),
		UserDisplayName:                         convert.GraphToFrameworkString(data.GetUserDisplayName()),
		WiFiMacAddress:                          convert.GraphToFrameworkString(data.GetWiFiMacAddress()),
		SubscriberCarrier:                       convert.GraphToFrameworkString(data.GetSubscriberCarrier()),
		MEID:                                    convert.GraphToFrameworkString(data.GetMeid()),
		TotalStorageSpaceInBytes:                convert.GraphToFrameworkInt64(data.GetTotalStorageSpaceInBytes()),
		FreeStorageSpaceInBytes:                 convert.GraphToFrameworkInt64(data.GetFreeStorageSpaceInBytes()),
		ManagedDeviceName:                       convert.GraphToFrameworkString(data.GetManagedDeviceName()),
		PartnerReportedThreatState:              convert.GraphToFrameworkEnum(data.GetPartnerReportedThreatState()),
		RetireAfterDateTime:                     convert.GraphToFrameworkTime(data.GetRetireAfterDateTime()),
		PreferMdmOverGroupPolicyAppliedDateTime: convert.GraphToFrameworkTime(data.GetPreferMdmOverGroupPolicyAppliedDateTime()),
		AutopilotEnrolled:                       convert.GraphToFrameworkBool(data.GetAutopilotEnrolled()),
		RequireUserEnrollmentApproval:           convert.GraphToFrameworkBool(data.GetRequireUserEnrollmentApproval()),
		ManagementCertificateExpirationDate:     convert.GraphToFrameworkTime(data.GetManagementCertificateExpirationDate()),
		ICCID:                                   convert.GraphToFrameworkString(data.GetIccid()),
		UDID:                                    convert.GraphToFrameworkString(data.GetUdid()),
		Notes:                                   convert.GraphToFrameworkString(data.GetNotes()),
		EthernetMacAddress:                      convert.GraphToFrameworkString(data.GetEthernetMacAddress()),
		PhysicalMemoryInBytes:                   convert.GraphToFrameworkInt64(data.GetPhysicalMemoryInBytes()),
		ProcessorArchitecture:                   convert.GraphToFrameworkEnum(data.GetProcessorArchitecture()),
		SpecificationVersion:                    convert.GraphToFrameworkString(data.GetSpecificationVersion()),
		JoinType:                                convert.GraphToFrameworkEnum(data.GetJoinType()),
		SkuFamily:                               convert.GraphToFrameworkString(data.GetSkuFamily()),
		SecurityPatchLevel:                      convert.GraphToFrameworkString(data.GetSecurityPatchLevel()),
		SkuNumber:                               convert.GraphToFrameworkInt32(data.GetSkuNumber()),
		ManagementFeatures:                      convert.GraphToFrameworkEnum(data.GetManagementFeatures()),
		EnrollmentProfileName:                   convert.GraphToFrameworkString(data.GetEnrollmentProfileName()),
		BootstrapTokenEscrowed:                  convert.GraphToFrameworkBool(data.GetBootstrapTokenEscrowed()),
		DeviceFirmwareConfigurationInterfaceManaged: convert.GraphToFrameworkBool(data.GetDeviceFirmwareConfigurationInterfaceManaged()),
		HardwareInformation:                         mapHardwareInformation(data.GetHardwareInformation()),
		DeviceActionResults:                         mapDeviceActionResults(data.GetDeviceActionResults()),
		ConfigurationManagerClientEnabledFeatures:   mapConfigurationManagerClientEnabledFeatures(data.GetConfigurationManagerClientEnabledFeatures()),
		DeviceHealthAttestationState:                mapDeviceHealthAttestationState(data.GetDeviceHealthAttestationState()),
		UsersLoggedOn:                               mapUsersLoggedOn(data.GetUsersLoggedOn()),
		RoleScopeTagIds:                             convert.GraphToFrameworkStringSlice(data.GetRoleScopeTagIds()),
		WindowsActiveMalwareCount:                   convert.GraphToFrameworkInt32AsInt64(data.GetWindowsActiveMalwareCount()),
		WindowsRemediatedMalwareCount:               convert.GraphToFrameworkInt32AsInt64(data.GetWindowsRemediatedMalwareCount()),
		ConfigurationManagerClientHealthState:       mapConfigurationManagerClientHealthState(data.GetConfigurationManagerClientHealthState()),
		ConfigurationManagerClientInformation:       mapConfigurationManagerClientInformation(data.GetConfigurationManagerClientInformation()),
		ChromeOSDeviceInfo:                          mapChromeOSDeviceInfo(data.GetChromeOSDeviceInfo()),
	}
}

// --- Mapping helpers for nested/complex fields ---

func mapHardwareInformation(hw graphmodels.HardwareInformationable) *ManagedDeviceHardwareInformation {
	if hw == nil {
		return nil
	}
	return &ManagedDeviceHardwareInformation{
		SerialNumber:                  convert.GraphToFrameworkString(hw.GetSerialNumber()),
		TotalStorageSpace:             convert.GraphToFrameworkInt64(hw.GetTotalStorageSpace()),
		FreeStorageSpace:              convert.GraphToFrameworkInt64(hw.GetFreeStorageSpace()),
		IMEI:                          convert.GraphToFrameworkString(hw.GetImei()),
		MEID:                          convert.GraphToFrameworkString(hw.GetMeid()),
		Manufacturer:                  convert.GraphToFrameworkString(hw.GetManufacturer()),
		Model:                         convert.GraphToFrameworkString(hw.GetModel()),
		PhoneNumber:                   convert.GraphToFrameworkString(hw.GetPhoneNumber()),
		SubscriberCarrier:             convert.GraphToFrameworkString(hw.GetSubscriberCarrier()),
		CellularTechnology:            convert.GraphToFrameworkString(hw.GetCellularTechnology()),
		WifiMac:                       convert.GraphToFrameworkString(hw.GetWifiMac()),
		OperatingSystemLanguage:       convert.GraphToFrameworkString(hw.GetOperatingSystemLanguage()),
		IsSupervised:                  convert.GraphToFrameworkBool(hw.GetIsSupervised()),
		IsEncrypted:                   convert.GraphToFrameworkBool(hw.GetIsEncrypted()),
		BatterySerialNumber:           convert.GraphToFrameworkString(hw.GetBatterySerialNumber()),
		BatteryHealthPercentage:       convert.GraphToFrameworkInt32AsInt64(hw.GetBatteryHealthPercentage()),
		BatteryChargeCycles:           convert.GraphToFrameworkInt32AsInt64(hw.GetBatteryChargeCycles()),
		IsSharedDevice:                convert.GraphToFrameworkBool(hw.GetIsSharedDevice()),
		SharedDeviceCachedUsers:       mapSharedAppleDeviceUsers(hw.GetSharedDeviceCachedUsers()),
		TPMSpecificationVersion:       convert.GraphToFrameworkString(hw.GetTpmSpecificationVersion()),
		OperatingSystemEdition:        convert.GraphToFrameworkString(hw.GetOperatingSystemEdition()),
		DeviceFullQualifiedDomainName: convert.GraphToFrameworkString(hw.GetDeviceFullQualifiedDomainName()),
		DeviceGuardVirtualizationBasedSecurityHardwareRequirementState: convert.GraphToFrameworkEnum(hw.GetDeviceGuardVirtualizationBasedSecurityHardwareRequirementState()),
		DeviceGuardVirtualizationBasedSecurityState:                    convert.GraphToFrameworkEnum(hw.GetDeviceGuardVirtualizationBasedSecurityState()),
		DeviceGuardLocalSystemAuthorityCredentialGuardState:            convert.GraphToFrameworkEnum(hw.GetDeviceGuardLocalSystemAuthorityCredentialGuardState()),
		OSBuildNumber:                       convert.GraphToFrameworkString(hw.GetOsBuildNumber()),
		OperatingSystemProductType:          convert.GraphToFrameworkInt32AsInt64(hw.GetOperatingSystemProductType()),
		IPAddressV4:                         convert.GraphToFrameworkString(hw.GetIpAddressV4()),
		SubnetAddress:                       convert.GraphToFrameworkString(hw.GetSubnetAddress()),
		ESIMIdentifier:                      convert.GraphToFrameworkString(hw.GetEsimIdentifier()),
		SystemManagementBIOSVersion:         convert.GraphToFrameworkString(hw.GetSystemManagementBIOSVersion()),
		TPMManufacturer:                     convert.GraphToFrameworkString(hw.GetTpmManufacturer()),
		TPMVersion:                          convert.GraphToFrameworkString(hw.GetTpmVersion()),
		WiredIPv4Addresses:                  convert.GraphToFrameworkStringSlice(hw.GetWiredIPv4Addresses()),
		BatteryLevelPercentage:              convert.GraphToFrameworkFloat64(hw.GetBatteryLevelPercentage()),
		ResidentUsersCount:                  convert.GraphToFrameworkInt32AsInt64(hw.GetResidentUsersCount()),
		ProductName:                         convert.GraphToFrameworkString(hw.GetProductName()),
		DeviceLicensingStatus:               convert.GraphToFrameworkEnum(hw.GetDeviceLicensingStatus()),
		DeviceLicensingLastErrorCode:        convert.GraphToFrameworkInt32AsInt64(hw.GetDeviceLicensingLastErrorCode()),
		DeviceLicensingLastErrorDescription: convert.GraphToFrameworkString(hw.GetDeviceLicensingLastErrorDescription()),
	}
}

func mapSharedAppleDeviceUsers(users []graphmodels.SharedAppleDeviceUserable) []SharedAppleDeviceUser {
	var result []SharedAppleDeviceUser
	for _, u := range users {
		result = append(result, SharedAppleDeviceUser{
			UserPrincipalName: convert.GraphToFrameworkString(u.GetUserPrincipalName()),
			DataToSync:        convert.GraphToFrameworkBool(u.GetDataToSync()),
			DataQuota:         convert.GraphToFrameworkInt64(u.GetDataQuota()),
			DataUsed:          convert.GraphToFrameworkInt64(u.GetDataUsed()),
		})
	}
	return result
}

func mapDeviceActionResults(results []graphmodels.DeviceActionResultable) []DeviceActionResult {
	var out []DeviceActionResult
	for _, r := range results {
		out = append(out, DeviceActionResult{
			ActionName:          convert.GraphToFrameworkString(r.GetActionName()),
			ActionState:         convert.GraphToFrameworkEnum(r.GetActionState()),
			StartDateTime:       convert.GraphToFrameworkTime(r.GetStartDateTime()),
			LastUpdatedDateTime: convert.GraphToFrameworkTime(r.GetLastUpdatedDateTime()),
		})
	}
	return out
}

func mapUsersLoggedOn(users []graphmodels.LoggedOnUserable) []LoggedOnUser {
	var out []LoggedOnUser
	for _, u := range users {
		out = append(out, LoggedOnUser{
			UserId:            convert.GraphToFrameworkString(u.GetUserId()),
			LastLogOnDateTime: convert.GraphToFrameworkTime(u.GetLastLogOnDateTime()),
		})
	}
	return out
}

func mapConfigurationManagerClientEnabledFeatures(src graphmodels.ConfigurationManagerClientEnabledFeaturesable) *ConfigurationManagerClientEnabledFeatures {
	if src == nil {
		return nil
	}
	return &ConfigurationManagerClientEnabledFeatures{
		Inventory:                convert.GraphToFrameworkBool(src.GetInventory()),
		ModernApps:               convert.GraphToFrameworkBool(src.GetModernApps()),
		ResourceAccess:           convert.GraphToFrameworkBool(src.GetResourceAccess()),
		DeviceConfiguration:      convert.GraphToFrameworkBool(src.GetDeviceConfiguration()),
		CompliancePolicy:         convert.GraphToFrameworkBool(src.GetCompliancePolicy()),
		WindowsUpdateForBusiness: convert.GraphToFrameworkBool(src.GetWindowsUpdateForBusiness()),
		EndpointProtection:       convert.GraphToFrameworkBool(src.GetEndpointProtection()),
		OfficeApps:               convert.GraphToFrameworkBool(src.GetOfficeApps()),
	}
}

func mapDeviceHealthAttestationState(data graphmodels.DeviceHealthAttestationStateable) *DeviceHealthAttestationState {
	if data == nil {
		return nil
	}
	return &DeviceHealthAttestationState{
		AttestationIdentityKey:                   convert.GraphToFrameworkString(data.GetAttestationIdentityKey()),
		BitLockerStatus:                          convert.GraphToFrameworkString(data.GetBitLockerStatus()),
		BootAppSecurityVersion:                   convert.GraphToFrameworkString(data.GetBootAppSecurityVersion()),
		BootDebugging:                            convert.GraphToFrameworkString(data.GetBootDebugging()),
		BootManagerSecurityVersion:               convert.GraphToFrameworkString(data.GetBootManagerSecurityVersion()),
		BootManagerVersion:                       convert.GraphToFrameworkString(data.GetBootManagerVersion()),
		BootRevisionListInfo:                     convert.GraphToFrameworkString(data.GetBootRevisionListInfo()),
		CodeIntegrity:                            convert.GraphToFrameworkString(data.GetCodeIntegrity()),
		CodeIntegrityCheckVersion:                convert.GraphToFrameworkString(data.GetCodeIntegrityCheckVersion()),
		CodeIntegrityPolicy:                      convert.GraphToFrameworkString(data.GetCodeIntegrityPolicy()),
		ContentNamespaceUrl:                      convert.GraphToFrameworkString(data.GetContentNamespaceUrl()),
		ContentVersion:                           convert.GraphToFrameworkString(data.GetContentVersion()),
		DataExcutionPolicy:                       convert.GraphToFrameworkString(data.GetDataExcutionPolicy()),
		DeviceHealthAttestationStatus:            convert.GraphToFrameworkString(data.GetDeviceHealthAttestationStatus()),
		EarlyLaunchAntiMalwareDriverProtection:   convert.GraphToFrameworkString(data.GetEarlyLaunchAntiMalwareDriverProtection()),
		FirmwareProtection:                       convert.GraphToFrameworkEnum(data.GetFirmwareProtection()),
		HealthAttestationSupportedStatus:         convert.GraphToFrameworkString(data.GetHealthAttestationSupportedStatus()),
		HealthStatusMismatchInfo:                 convert.GraphToFrameworkString(data.GetHealthStatusMismatchInfo()),
		IssuedDateTime:                           convert.GraphToFrameworkTime(data.GetIssuedDateTime()),
		LastUpdateDateTime:                       convert.GraphToFrameworkString(data.GetLastUpdateDateTime()),
		MemoryAccessProtection:                   convert.GraphToFrameworkEnum(data.GetMemoryAccessProtection()),
		MemoryIntegrityProtection:                convert.GraphToFrameworkEnum(data.GetMemoryIntegrityProtection()),
		OperatingSystemKernelDebugging:           convert.GraphToFrameworkString(data.GetOperatingSystemKernelDebugging()),
		OperatingSystemRevListInfo:               convert.GraphToFrameworkString(data.GetOperatingSystemRevListInfo()),
		Pcr0:                                     convert.GraphToFrameworkString(data.GetPcr0()),
		PcrHashAlgorithm:                         convert.GraphToFrameworkString(data.GetPcrHashAlgorithm()),
		ResetCount:                               convert.GraphToFrameworkInt64(data.GetResetCount()),
		RestartCount:                             convert.GraphToFrameworkInt64(data.GetRestartCount()),
		SafeMode:                                 convert.GraphToFrameworkString(data.GetSafeMode()),
		SecureBoot:                               convert.GraphToFrameworkString(data.GetSecureBoot()),
		SecureBootConfigurationPolicyFingerPrint: convert.GraphToFrameworkString(data.GetSecureBootConfigurationPolicyFingerPrint()),
		SecuredCorePC:                            convert.GraphToFrameworkEnum(data.GetSecuredCorePC()),
		SystemManagementMode:                     convert.GraphToFrameworkEnum(data.GetSystemManagementMode()),
		TestSigning:                              convert.GraphToFrameworkString(data.GetTestSigning()),
		TpmVersion:                               convert.GraphToFrameworkString(data.GetTpmVersion()),
		VirtualizationBasedSecurity:              convert.GraphToFrameworkEnum(data.GetVirtualizationBasedSecurity()),
		VirtualSecureMode:                        convert.GraphToFrameworkString(data.GetVirtualSecureMode()),
		WindowsPE:                                convert.GraphToFrameworkString(data.GetWindowsPE()),
	}
}

func mapConfigurationManagerClientHealthState(src graphmodels.ConfigurationManagerClientHealthStateable) *ConfigurationManagerClientHealthState {
	if src == nil {
		return nil
	}
	return &ConfigurationManagerClientHealthState{
		State:            convert.GraphToFrameworkEnum(src.GetState()),
		ErrorCode:        convert.GraphToFrameworkInt32AsInt64(src.GetErrorCode()),
		LastSyncDateTime: convert.GraphToFrameworkTime(src.GetLastSyncDateTime()),
	}
}

func mapConfigurationManagerClientInformation(src graphmodels.ConfigurationManagerClientInformationable) *ConfigurationManagerClientInformation {
	if src == nil {
		return nil
	}
	return &ConfigurationManagerClientInformation{
		ClientIdentifier: convert.GraphToFrameworkString(src.GetClientIdentifier()),
		IsBlocked:        convert.GraphToFrameworkBool(src.GetIsBlocked()),
		ClientVersion:    convert.GraphToFrameworkString(src.GetClientVersion()),
	}
}

func mapChromeOSDeviceInfo(src []graphmodels.ChromeOSDevicePropertyable) []ChromeOSDeviceInfo {
	var out []ChromeOSDeviceInfo
	for _, c := range src {
		out = append(out, ChromeOSDeviceInfo{
			Name:      convert.GraphToFrameworkString(c.GetName()),
			Value:     convert.GraphToFrameworkString(c.GetValue()),
			ValueType: convert.GraphToFrameworkString(c.GetValueType()),
			Updatable: convert.GraphToFrameworkBool(c.GetUpdatable()),
		})
	}
	return out
}
