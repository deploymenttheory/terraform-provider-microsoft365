package graphBetaRoleDefinition

// This file contains constants for all Microsoft Intune role permissions
// These constants are used when constructing role definitions

// Common prefix for all Intune permissions
const IntunePermissionPrefix = "Microsoft.Intune_"

// Audit permissions
const (
	IntuneAuditRead = "Microsoft.Intune_Audit_Read"
)

// Android related permissions
const (
	IntuneAndroidSyncRead                           = "Microsoft.Intune_AndroidSync_Read"
	IntuneAndroidSyncUpdateApps                     = "Microsoft.Intune_AndroidSync_UpdateApps"
	IntuneAndroidSyncUpdateOnboarding               = "Microsoft.Intune_AndroidSync_UpdateOnboarding"
	IntuneAndroidSyncUpdateEnrollmentProfiles       = "Microsoft.Intune_AndroidSync_UpdateEnrollmentProfiles"
	IntuneAndroidSyncAndroidEnrollmentTimeMemAssign = "Microsoft.Intune_AndroidSync_AndroidEnrollmentTimeMembershipAssign"
)

// Android FOTA (Firmware Over The Air) permissions
const (
	IntuneAndroidFotaRead   = "Microsoft.Intune_AndroidFota_Read"
	IntuneAndroidFotaCreate = "Microsoft.Intune_AndroidFota_Create"
	IntuneAndroidFotaUpdate = "Microsoft.Intune_AndroidFota_Update"
	IntuneAndroidFotaDelete = "Microsoft.Intune_AndroidFota_Delete"
	IntuneAndroidFotaAssign = "Microsoft.Intune_AndroidFota_Assign"
)

// Apple related permissions
const (
	IntuneAppleDeviceSerialNumbersRead   = "Microsoft.Intune_AppleDeviceSerialNumbers_Read"
	IntuneAppleDeviceSerialNumbersCreate = "Microsoft.Intune_AppleDeviceSerialNumbers_Create"
	IntuneAppleDeviceSerialNumbersUpdate = "Microsoft.Intune_AppleDeviceSerialNumbers_Update"
	IntuneAppleDeviceSerialNumbersDelete = "Microsoft.Intune_AppleDeviceSerialNumbers_Delete"
)

// Apple Enrollment Profiles permissions
const (
	IntuneAppleEnrollmentProfilesRead   = "Microsoft.Intune_AppleEnrollmentProfiles_Read"
	IntuneAppleEnrollmentProfilesCreate = "Microsoft.Intune_AppleEnrollmentProfiles_Create"
	IntuneAppleEnrollmentProfilesUpdate = "Microsoft.Intune_AppleEnrollmentProfiles_Update"
	IntuneAppleEnrollmentProfilesDelete = "Microsoft.Intune_AppleEnrollmentProfiles_Delete"
	IntuneAppleEnrollmentProfilesAssign = "Microsoft.Intune_AppleEnrollmentProfiles_Assign"
)

// App Control Policy permissions
const (
	IntuneAppControlPolicyRead        = "Microsoft.Intune_AppControlPolicy_Read"
	IntuneAppControlPolicyCreate      = "Microsoft.Intune_AppControlPolicy_Create"
	IntuneAppControlPolicyUpdate      = "Microsoft.Intune_AppControlPolicy_Update"
	IntuneAppControlPolicyDelete      = "Microsoft.Intune_AppControlPolicy_Delete"
	IntuneAppControlPolicyAssign      = "Microsoft.Intune_AppControlPolicy_Assign"
	IntuneAppControlPolicyViewReports = "Microsoft.Intune_AppControlPolicy_ViewReports"
)

// ASR (Attack Surface Reduction) Policy permissions
const (
	IntuneASRPolicyRead        = "Microsoft.Intune_ASRPolicy_Read"
	IntuneASRPolicyCreate      = "Microsoft.Intune_ASRPolicy_Create"
	IntuneASRPolicyUpdate      = "Microsoft.Intune_ASRPolicy_Update"
	IntuneASRPolicyDelete      = "Microsoft.Intune_ASRPolicy_Delete"
	IntuneASRPolicyAssign      = "Microsoft.Intune_ASRPolicy_Assign"
	IntuneASRPolicyViewReports = "Microsoft.Intune_ASRPolicy_ViewReports"
)

// Assignment Filter permissions
const (
	IntuneAssignmentFilterRead   = "Microsoft.Intune_AssignmentFilter_Read"
	IntuneAssignmentFilterCreate = "Microsoft.Intune_AssignmentFilter_Create"
	IntuneAssignmentFilterUpdate = "Microsoft.Intune_AssignmentFilter_Update"
	IntuneAssignmentFilterDelete = "Microsoft.Intune_AssignmentFilter_Delete"
)

// Certificate Connector permissions
const (
	IntuneCertificateConnectorRead   = "Microsoft.Intune_CertificateConnector_Read"
	IntuneCertificateConnectorModify = "Microsoft.Intune_CertificateConnector_Modify"
)

// Chromebook Sync permissions
const (
	IntuneChromebookSyncRead             = "Microsoft.Intune_ChromebookSync_Read"
	IntuneChromebookSyncUpdateOnboarding = "Microsoft.Intune_ChromebookSync_UpdateOnboarding"
	IntuneChromebookSyncDeleteOnboarding = "Microsoft.Intune_ChromebookSync_DeleteOnboarding"
)

// CloudAttach permissions
const (
	IntuneCloudAttachResourceExplorer   = "Microsoft.Intune_CloudAttach_ResourceExplorer"
	IntuneCloudAttachTimeline           = "Microsoft.Intune_CloudAttach_Timeline"
	IntuneCloudAttachCollections        = "Microsoft.Intune_CloudAttach_Collections"
	IntuneCloudAttachScripts            = "Microsoft.Intune_CloudAttach_Scripts"
	IntuneCloudAttachScriptActions      = "Microsoft.Intune_CloudAttach_ScriptActions"
	IntuneCloudAttachSoftwareUpdates    = "Microsoft.Intune_CloudAttach_SoftwareUpdates"
	IntuneCloudAttachCMPivot            = "Microsoft.Intune_CloudAttach_CMPivot"
	IntuneCloudAttachApplicationActions = "Microsoft.Intune_CloudAttach_ApplicationActions"
	IntuneCloudAttachClientDetails      = "Microsoft.Intune_CloudAttach_ClientDetails"
	IntuneCloudAttachApplications       = "Microsoft.Intune_CloudAttach_Applications"
	IntuneCloudAttachEnrollNow          = "Microsoft.Intune_CloudAttach_EnrollNow"
)

// CloudPki permissions
const (
	IntuneCloudPkiRead   = "Microsoft.Intune_CloudPki_Read"
	IntuneCloudPkiCreate = "Microsoft.Intune_CloudPki_Create"
	IntuneCloudPkiUpdate = "Microsoft.Intune_CloudPki_Update"
	IntuneCloudPkiRevoke = "Microsoft.Intune_CloudPki_Revoke"
)

// Corporate Device Identifiers permissions
const (
	IntuneCorporateDeviceIdentifiersRead   = "Microsoft.Intune_CorporateDeviceIdentifiers_Read"
	IntuneCorporateDeviceIdentifiersCreate = "Microsoft.Intune_CorporateDeviceIdentifiers_Create"
	IntuneCorporateDeviceIdentifiersUpdate = "Microsoft.Intune_CorporateDeviceIdentifiers_Update"
	IntuneCorporateDeviceIdentifiersDelete = "Microsoft.Intune_CorporateDeviceIdentifiers_Delete"
)

// Customization permissions
const (
	IntuneCustomizationRead   = "Microsoft.Intune_Customization_Read"
	IntuneCustomizationCreate = "Microsoft.Intune_Customization_Create"
	IntuneCustomizationUpdate = "Microsoft.Intune_Customization_Update"
	IntuneCustomizationDelete = "Microsoft.Intune_Customization_Delete"
	IntuneCustomizationAssign = "Microsoft.Intune_Customization_Assign"
)

// Derived Credentials permissions
const (
	IntuneDerivedCredentialsRead   = "Microsoft.Intune_DerivedCredentials_Read"
	IntuneDerivedCredentialsModify = "Microsoft.Intune_DerivedCredentials_Modify"
)

// Device Compliance Policies permissions
const (
	IntuneDeviceCompliancePoliciesRead        = "Microsoft.Intune_DeviceCompliancePolices_Read"
	IntuneDeviceCompliancePoliciesCreate      = "Microsoft.Intune_DeviceCompliancePolices_Create"
	IntuneDeviceCompliancePoliciesUpdate      = "Microsoft.Intune_DeviceCompliancePolices_Update"
	IntuneDeviceCompliancePoliciesDelete      = "Microsoft.Intune_DeviceCompliancePolices_Delete"
	IntuneDeviceCompliancePoliciesAssign      = "Microsoft.Intune_DeviceCompliancePolices_Assign"
	IntuneDeviceCompliancePoliciesViewReports = "Microsoft.Intune_DeviceCompliancePolices_ViewReports"
)

// Device Configurations permissions
const (
	IntuneDeviceConfigurationsRead                 = "Microsoft.Intune_DeviceConfigurations_Read"
	IntuneDeviceConfigurationsCreate               = "Microsoft.Intune_DeviceConfigurations_Create"
	IntuneDeviceConfigurationsUpdate               = "Microsoft.Intune_DeviceConfigurations_Update"
	IntuneDeviceConfigurationsDelete               = "Microsoft.Intune_DeviceConfigurations_Delete"
	IntuneDeviceConfigurationsAssign               = "Microsoft.Intune_DeviceConfigurations_Assign"
	IntuneDeviceConfigurationsViewReports          = "Microsoft.Intune_DeviceConfigurations_ViewReports"
	IntuneDeviceConfigurationsUpdateWindowsRestore = "Microsoft.Intune_DeviceConfigurations_UpdateWindowsRestore"
)

// Device Enrollment Managers permissions
const (
	IntuneDeviceEnrollmentManagersRead   = "Microsoft.Intune_DeviceEnrollmentManagers_Read"
	IntuneDeviceEnrollmentManagersUpdate = "Microsoft.Intune_DeviceEnrollmentManagers_Update"
)

// EDR (Endpoint Detection and Response) Policy permissions
const (
	IntuneEDRPolicyRead        = "Microsoft.Intune_EDRPolicy_Read"
	IntuneEDRPolicyCreate      = "Microsoft.Intune_EDRPolicy_Create"
	IntuneEDRPolicyUpdate      = "Microsoft.Intune_EDRPolicy_Update"
	IntuneEDRPolicyDelete      = "Microsoft.Intune_EDRPolicy_Delete"
	IntuneEDRPolicyAssign      = "Microsoft.Intune_EDRPolicy_Assign"
	IntuneEDRPolicyViewReports = "Microsoft.Intune_EDRPolicy_ViewReports"
)

// Endpoint Analytics permissions
const (
	IntuneEndpointAnalyticsRead   = "Microsoft.Intune_EndpointAnalytics_Read"
	IntuneEndpointAnalyticsCreate = "Microsoft.Intune_EndpointAnalytics_Create"
	IntuneEndpointAnalyticsUpdate = "Microsoft.Intune_EndpointAnalytics_Update"
	IntuneEndpointAnalyticsDelete = "Microsoft.Intune_EndpointAnalytics_Delete"
)

// Endpoint Protection permissions
const (
	IntuneEndpointProtectionRead = "Microsoft.Intune_EndpointProtection_Read"
)

// Enrollment Program Token permissions
const (
	IntuneEnrollmentProgramTokenRead   = "Microsoft.Intune_EnrollmentProgramToken_Read"
	IntuneEnrollmentProgramTokenCreate = "Microsoft.Intune_EnrollmentProgramToken_Create"
	IntuneEnrollmentProgramTokenUpdate = "Microsoft.Intune_EnrollmentProgramToken_Update"
	IntuneEnrollmentProgramTokenDelete = "Microsoft.Intune_EnrollmentProgramToken_Delete"
)

// Enrollment Profiles permissions
const (
	IntuneEnrollmentProfilesEnrollmentTimeMembershipAssign = "Microsoft.Intune_EnrollmentProfiles_EnrollmentTimeMembershipAssign"
)

// EPM (Endpoint Privilege Management) Policy permissions
const (
	IntuneEpmPolicyRead                    = "Microsoft.Intune_EpmPolicy_Read"
	IntuneEpmPolicyCreate                  = "Microsoft.Intune_EpmPolicy_Create"
	IntuneEpmPolicyUpdate                  = "Microsoft.Intune_EpmPolicy_Update"
	IntuneEpmPolicyDelete                  = "Microsoft.Intune_EpmPolicy_Delete"
	IntuneEpmPolicyAssign                  = "Microsoft.Intune_EpmPolicy_Assign"
	IntuneEpmPolicyViewReports             = "Microsoft.Intune_EpmPolicy_ViewReports"
	IntuneEpmPolicyViewElevationRequests   = "Microsoft.Intune_EpmPolicy_ViewElevationRequests"
	IntuneEpmPolicyModifyElevationRequests = "Microsoft.Intune_EpmPolicy_ModifyElevationRequests"
)

// Managed Apps permissions
const (
	IntuneManagedAppsRead   = "Microsoft.Intune_ManagedApps_Read"
	IntuneManagedAppsCreate = "Microsoft.Intune_ManagedApps_Create"
	IntuneManagedAppsUpdate = "Microsoft.Intune_ManagedApps_Update"
	IntuneManagedAppsDelete = "Microsoft.Intune_ManagedApps_Delete"
	IntuneManagedAppsAssign = "Microsoft.Intune_ManagedApps_Assign"
	IntuneManagedAppsWipe   = "Microsoft.Intune_ManagedApps_Wipe"
)

// Managed Device Cleanup permissions
const (
	IntuneManagedDeviceCleanupSettingsUpdate = "Microsoft.Intune_ManagedDeviceCleanupSettings_Update"
	IntuneManagedDeviceCleanupRuleUpdate     = "Microsoft.Intune_ManagedDeviceCleanupRule_Update"
)

// Managed Devices permissions
const (
	IntuneManagedDevicesRead             = "Microsoft.Intune_ManagedDevices_Read"
	IntuneManagedDevicesUpdate           = "Microsoft.Intune_ManagedDevices_Update"
	IntuneManagedDevicesDelete           = "Microsoft.Intune_ManagedDevices_Delete"
	IntuneManagedDevicesSetPrimaryUser   = "Microsoft.Intune_ManagedDevices_SetPrimaryUser"
	IntuneManagedDevicesViewReports      = "Microsoft.Intune_ManagedDevices_ViewReports"
	IntuneManagedDevicesQuery            = "Microsoft.Intune_ManagedDevices_Query"
	IntuneManagedDevicesReadBiosPassword = "Microsoft.Intune_ManagedDevices_ReadBiosPassword"
)

// Managed Google Play permissions
const (
	IntuneManagedGooglePlayRead   = "Microsoft.Intune_ManagedGooglePlay_Read"
	IntuneManagedGooglePlayModify = "Microsoft.Intune_ManagedGooglePlay_Modify"
)

// Microsoft Defender ATP permissions
const (
	IntuneMicrosoftDefenderATPRead = "Microsoft.Intune_MicrosoftDefenderATP_Read"
)

// Microsoft Store for Business permissions
const (
	IntuneMicrosoftStoreForBusinessRead   = "Microsoft.Intune_MicrosoftStoreForBusiness_Read"
	IntuneMicrosoftStoreForBusinessModify = "Microsoft.Intune_MicrosoftStoreForBusiness_Modify"
)

// Microsoft Tunnel Gateway permissions
const (
	IntuneMicrosoftTunnelGatewayRead   = "Microsoft.Intune_MicrosoftTunnelGateway_Read"
	IntuneMicrosoftTunnelGatewayCreate = "Microsoft.Intune_MicrosoftTunnelGateway_Create"
	IntuneMicrosoftTunnelGatewayUpdate = "Microsoft.Intune_MicrosoftTunnelGateway_Update"
	IntuneMicrosoftTunnelGatewayDelete = "Microsoft.Intune_MicrosoftTunnelGateway_Delete"
)

// Mobile Apps permissions
const (
	IntuneMobileAppsRead        = "Microsoft.Intune_MobileApps_Read"
	IntuneMobileAppsCreate      = "Microsoft.Intune_MobileApps_Create"
	IntuneMobileAppsUpdate      = "Microsoft.Intune_MobileApps_Update"
	IntuneMobileAppsDelete      = "Microsoft.Intune_MobileApps_Delete"
	IntuneMobileAppsAssign      = "Microsoft.Intune_MobileApps_Assign"
	IntuneMobileAppsViewReports = "Microsoft.Intune_MobileApps_ViewReports"
	IntuneMobileAppsRelate      = "Microsoft.Intune_MobileApps_Relate"
)

// Mobile Threat Defense permissions
const (
	IntuneMobileThreatDefenseRead   = "Microsoft.Intune_MobileThreatDefense_Read"
	IntuneMobileThreatDefenseModify = "Microsoft.Intune_MobileThreatDefense_Modify"
)

// Multi-Admin Approval permissions
const (
	IntuneMultiAdminApprovalReadAccessPolicy      = "Microsoft.Intune_MultiAdminApproval_ReadAccessPolicy"
	IntuneMultiAdminApprovalCreateAccessPolicy    = "Microsoft.Intune_MultiAdminApproval_CreateAccessPolicy"
	IntuneMultiAdminApprovalUpdateAccessPolicy    = "Microsoft.Intune_MultiAdminApproval_UpdateAccessPolicy"
	IntuneMultiAdminApprovalDeleteAccessPolicy    = "Microsoft.Intune_MultiAdminApproval_DeleteAccessPolicy"
	IntuneMultiAdminApprovalApprovalForMultiAdmin = "Microsoft.Intune_MultiAdminApproval_ApprovalForMultiAdminApproval"
)

// Organization permissions
const (
	IntuneOrganizationRead   = "Microsoft.Intune_Organization_Read"
	IntuneOrganizationCreate = "Microsoft.Intune_Organization_Create"
	IntuneOrganizationUpdate = "Microsoft.Intune_Organization_Update"
	IntuneOrganizationDelete = "Microsoft.Intune_Organization_Delete"
)

// Organizational Messages permissions
const (
	IntuneOrganizationalMessagesRead          = "Microsoft.Intune_OrganizationalMessages_Read"
	IntuneOrganizationalMessagesCreate        = "Microsoft.Intune_OrganizationalMessages_Create"
	IntuneOrganizationalMessagesUpdate        = "Microsoft.Intune_OrganizationalMessages_Update"
	IntuneOrganizationalMessagesDelete        = "Microsoft.Intune_OrganizationalMessages_Delete"
	IntuneOrganizationalMessagesAssign        = "Microsoft.Intune_OrganizationalMessages_Assign"
	IntuneOrganizationalMessagesUpdateControl = "Microsoft.Intune_OrganizationalMessages_UpdateControl"
)

// Partner Device Management permissions
const (
	IntunePartnerDeviceManagementRead   = "Microsoft.Intune_PartnerDeviceManagement_Read"
	IntunePartnerDeviceManagementModify = "Microsoft.Intune_PartnerDeviceManagement_Modify"
)

// Policy Sets permissions
const (
	IntunePolicySetsRead   = "Microsoft.Intune_PolicySets_Read"
	IntunePolicySetsCreate = "Microsoft.Intune_PolicySets_Create"
	IntunePolicySetsUpdate = "Microsoft.Intune_PolicySets_Update"
	IntunePolicySetsDelete = "Microsoft.Intune_PolicySets_Delete"
	IntunePolicySetsAssign = "Microsoft.Intune_PolicySets_Assign"
)

// Quiet Time Policies permissions
const (
	IntuneQuietTimePoliciesRead        = "Microsoft.Intune_QuietTimePolicies_Read"
	IntuneQuietTimePoliciesCreate      = "Microsoft.Intune_QuietTimePolicies_Create"
	IntuneQuietTimePoliciesUpdate      = "Microsoft.Intune_QuietTimePolicies_Update"
	IntuneQuietTimePoliciesDelete      = "Microsoft.Intune_QuietTimePolicies_Delete"
	IntuneQuietTimePoliciesAssign      = "Microsoft.Intune_QuietTimePolicies_Assign"
	IntuneQuietTimePoliciesViewReports = "Microsoft.Intune_QuietTimePolicies_ViewReports"
)

// Remote Assistance permissions
const (
	IntuneRemoteAssistanceRead        = "Microsoft.Intune_RemoteAssistance_Read"
	IntuneRemoteAssistanceUpdate      = "Microsoft.Intune_RemoteAssistance_Update"
	IntuneRemoteAssistanceViewReports = "Microsoft.Intune_RemoteAssistance_ViewReports"
)

// Remote Assistance App permissions
const (
	IntuneRemoteAssistanceAppViewScreen      = "Microsoft.Intune_RemoteAssistanceApp_ViewScreen"
	IntuneRemoteAssistanceAppElevation       = "Microsoft.Intune_RemoteAssistanceApp_Elevation"
	IntuneRemoteAssistanceAppTakeFullControl = "Microsoft.Intune_RemoteAssistanceApp_TakeFullControl"
	IntuneRemoteAssistanceAppUnattended      = "Microsoft.Intune_RemoteAssistanceApp_Unattended"
)

// Remote Tasks permissions
const (
	IntuneRemoteTasksActivateDeviceEsim           = "Microsoft.Intune_RemoteTasks_ActivateDeviceEsim"
	IntuneRemoteTasksCustomNotification           = "Microsoft.Intune_RemoteTasks_CustomNotification"
	IntuneRemoteTasksPlayLostModeSound            = "Microsoft.Intune_RemoteTasks_PlayLostModeSound"
	IntuneRemoteTasksDeviceLogs                   = "Microsoft.Intune_RemoteTasks_DeviceLogs"
	IntuneRemoteTasksRebootNow                    = "Microsoft.Intune_RemoteTasks_RebootNow"
	IntuneRemoteTasksRemoteLock                   = "Microsoft.Intune_RemoteTasks_RemoteLock"
	IntuneRemoteTasksRevokeAppleVppLicenses       = "Microsoft.Intune_RemoteTasks_RevokeAppleVppLicenses"
	IntuneRemoteTasksRotateFileVaultKey           = "Microsoft.Intune_RemoteTasks_RotateFileVaultKey"
	IntuneRemoteTasksSyncDevice                   = "Microsoft.Intune_RemoteTasks_SyncDevice"
	IntuneRemoteTasksEnableWindowsIntuneAgent     = "Microsoft.Intune_RemoteTasks_EnableWindowsIntuneAgent"
	IntuneRemoteTasksLocateDevice                 = "Microsoft.Intune_RemoteTasks_LocateDevice"
	IntuneRemoteTasksConfigurationManagerAction   = "Microsoft.Intune_RemoteTasks_ConfigurationManagerAction"
	IntuneRemoteTasksResetPasscode                = "Microsoft.Intune_RemoteTasks_ResetPasscode"
	IntuneRemoteTasksEnableLostMode               = "Microsoft.Intune_RemoteTasks_EnableLostMode"
	IntuneRemoteTasksDisableLostMode              = "Microsoft.Intune_RemoteTasks_DisableLostMode"
	IntuneRemoteTasksRotateBitLockerKeys          = "Microsoft.Intune_RemoteTasks_RotateBitLockerKeys"
	IntuneRemoteTasksGetFileVaultKey              = "Microsoft.Intune_RemoteTasks_GetFileVaultKey"
	IntuneRemoteTasksWindowsDefender              = "Microsoft.Intune_RemoteTasks_WindowsDefender"
	IntuneRemoteTasksRemoveDFCIManagement         = "Microsoft.Intune_RemoteTasks_RemoveDFCIManagement"
	IntuneRemoteTasksInitiateDeviceAttestation    = "Microsoft.Intune_RemoteTasks_InitiateDeviceAttestation"
	IntuneRemoteTasksCleanPC                      = "Microsoft.Intune_RemoteTasks_CleanPC"
	IntuneRemoteTasksShutDown                     = "Microsoft.Intune_RemoteTasks_ShutDown"
	IntuneRemoteTasksOnDemandProactiveRemediation = "Microsoft.Intune_RemoteTasks_OnDemandProactiveRemediation"
	IntuneRemoteTasksPauseConfigurationRefresh    = "Microsoft.Intune_RemoteTasks_PauseConfigurationRefresh"
	IntuneRemoteTasksManageSharedDeviceUsers      = "Microsoft.Intune_RemoteTasks_ManageSharedDeviceUsers"
	IntuneRemoteTasksRequestRemoteAssistance      = "Microsoft.Intune_RemoteTasks_RequestRemoteAssistance"
	IntuneRemoteTasksRetire                       = "Microsoft.Intune_RemoteTasks_Retire"
	IntuneRemoteTasksBypassActivationLock         = "Microsoft.Intune_RemoteTasks_BypassActivationLock"
	IntuneRemoteTasksUpdateDeviceAccount          = "Microsoft.Intune_RemoteTasks_UpdateDeviceAccount"
	IntuneRemoteTasksSetDeviceName                = "Microsoft.Intune_RemoteTasks_SetDeviceName"
	IntuneRemoteTasksWipe                         = "Microsoft.Intune_RemoteTasks_Wipe"
	IntuneRemoteTasksInitiateMDMKeyRecovery       = "Microsoft.Intune_RemoteTasks_InitiateMDMKeyRecovery"
	IntuneRemoteTasksRotateLocalAdminPassword     = "Microsoft.Intune_RemoteTasks_RotateLocalAdminPassword"
	IntuneRemoteTasksChangeAssignments            = "Microsoft.Intune_RemoteTasks_ChangeAssignments"
	IntuneRemoteTasksOffboard                     = "Microsoft.Intune_RemoteTasks_Offboard"
)

// Reports permissions
const (
	IntuneReportsRead = "Microsoft.Intune_Reports_Read"
)

// Roles permissions
const (
	IntuneRolesRead   = "Microsoft.Intune_Roles_Read"
	IntuneRolesCreate = "Microsoft.Intune_Roles_Create"
	IntuneRolesUpdate = "Microsoft.Intune_Roles_Update"
	IntuneRolesDelete = "Microsoft.Intune_Roles_Delete"
	IntuneRolesAssign = "Microsoft.Intune_Roles_Assign"
)

// Security Baselines permissions
const (
	IntuneSecurityBaselinesRead   = "Microsoft.Intune_SecurityBaselines_Read"
	IntuneSecurityBaselinesCreate = "Microsoft.Intune_SecurityBaselines_Create"
	IntuneSecurityBaselinesUpdate = "Microsoft.Intune_SecurityBaselines_Update"
	IntuneSecurityBaselinesDelete = "Microsoft.Intune_SecurityBaselines_Delete"
	IntuneSecurityBaselinesAssign = "Microsoft.Intune_SecurityBaselines_Assign"
)

// Security Tasks permissions
const (
	IntuneSecurityTasksRead   = "Microsoft.Intune_SecurityTasks_Read"
	IntuneSecurityTasksUpdate = "Microsoft.Intune_SecurityTasks_Update"
)

// ServiceNow permissions
const (
	IntuneServiceNowUpdateConnector = "Microsoft.Intune_ServiceNow_UpdateConnector"
	IntuneServiceNowViewIncidents   = "Microsoft.Intune_ServiceNow_ViewIncidents"
)

// Telecom Expenses permissions
const (
	IntuneTelecomExpensesRead   = "Microsoft.Intune_TelecomExpenses_Read"
	IntuneTelecomExpensesUpdate = "Microsoft.Intune_TelecomExpenses_Update"
)

// Tenant Attach Recommendations permissions
const (
	IntuneTenantAttachRecommendationsRead = "Microsoft.Intune_TenantAttachRecommendations_Read"
)

// Terms and Conditions permissions
const (
	IntuneTermsAndConditionsRead   = "Microsoft.Intune_TermsAndConditions_Read"
	IntuneTermsAndConditionsCreate = "Microsoft.Intune_TermsAndConditions_Create"
	IntuneTermsAndConditionsUpdate = "Microsoft.Intune_TermsAndConditions_Update"
	IntuneTermsAndConditionsDelete = "Microsoft.Intune_TermsAndConditions_Delete"
	IntuneTermsAndConditionsAssign = "Microsoft.Intune_TermsAndConditions_Assign"
)

// Windows Enterprise Certificate permissions
const (
	IntuneWindowsEnterpriseCertificateRead   = "Microsoft.Intune_WindowsEnterpriseCertificate_Read"
	IntuneWindowsEnterpriseCertificateModify = "Microsoft.Intune_WindowsEnterpriseCertificate_Modify"
)

// Windows OS Recovery permissions
const (
	IntuneWindowsOSRecoveryRead   = "Microsoft.Intune_WindowsOSRecovery_Read"
	IntuneWindowsOSRecoveryCreate = "Microsoft.Intune_WindowsOSRecovery_Create"
	IntuneWindowsOSRecoveryUpdate = "Microsoft.Intune_WindowsOSRecovery_Update"
	IntuneWindowsOSRecoveryDelete = "Microsoft.Intune_WindowsOSRecovery_Delete"
	IntuneWindowsOSRecoveryAssign = "Microsoft.Intune_WindowsOSRecovery_Assign"
)
