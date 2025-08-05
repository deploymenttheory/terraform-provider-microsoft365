Intune

POST https://graph.microsoft.com/beta/deviceManagement/roleDefinitions

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewRoleDefinition()
id := ""
requestBody.SetId(&id) 
description := "test"
requestBody.SetDescription(&description) 
displayName := "test"
requestBody.SetDisplayName(&displayName) 


rolePermission := graphmodels.NewRolePermission()


resourceAction := graphmodels.NewResourceAction()
allowedResourceActions := []string {
	"Microsoft.Intune_Audit_Read",
	"Microsoft.Intune_CorporateDeviceIdentifiers_Create",
	"Microsoft.Intune_EnrollmentProgramToken_Create",
	"Microsoft.Intune_ManagedDevices_Update",
	"Microsoft.Intune_AssignmentFilter_Read",
	"Microsoft.Intune_DeviceCompliancePolices_Assign",
	"Microsoft.Intune_MicrosoftTunnelGateway_Create",
	"Microsoft.Intune_ManagedDevices_Delete",
	"Microsoft.Intune_MicrosoftTunnelGateway_Delete",
	"Microsoft.Intune_AppleEnrollmentProfiles_Assign",
	"Microsoft.Intune_ManagedGooglePlay_Read",
	"Microsoft.Intune_CorporateDeviceIdentifiers_Update",
	"Microsoft.Intune_AndroidSync_UpdateApps",
	"Microsoft.Intune_DeviceCompliancePolices_Delete",
	"Microsoft.Intune_AppleEnrollmentProfiles_Delete",
	"Microsoft.Intune_MicrosoftStoreForBusiness_Read",
	"Microsoft.Intune_DeviceCompliancePolices_Update",
	"Microsoft.Intune_MicrosoftTunnelGateway_Update",
	"Microsoft.Intune_RemoteAssistanceApp_ViewScreen",
	"Microsoft.Intune_EnrollmentProgramToken_Update",
	"Microsoft.Intune_AppleEnrollmentProfiles_Create",
	"Microsoft.Intune_ManagedDevices_SetPrimaryUser",
	"Microsoft.Intune_AppleEnrollmentProfiles_Update",
	"Microsoft.Intune_RemoteTasks_ActivateDeviceEsim",
	"Microsoft.Intune_AppleDeviceSerialNumbers_Read",
	"Microsoft.Intune_EndpointAnalytics_Read",
	"Microsoft.Intune_CorporateDeviceIdentifiers_Delete",
	"Microsoft.Intune_RemoteTasks_CustomNotification",
	"Microsoft.Intune_RemoteTasks_PlayLostModeSound",
	"Microsoft.Intune_RemoteTasks_DeviceLogs",
	"Microsoft.Intune_RemoteTasks_RebootNow",
	"Microsoft.Intune_RemoteTasks_RemoteLock",
	"Microsoft.Intune_RemoteTasks_RevokeAppleVppLicenses",
	"Microsoft.Intune_RemoteTasks_RotateFileVaultKey",
	"Microsoft.Intune_RemoteTasks_SyncDevice",
	"Microsoft.Intune_AssignmentFilter_Create",
	"Microsoft.Intune_RemoteTasks_EnableWindowsIntuneAgent",
	"Microsoft.Intune_Roles_Create",
	"Microsoft.Intune_DerivedCredentials_Read",
	"Microsoft.Intune_EndpointAnalytics_Create",
	"Microsoft.Intune_Roles_Read",
	"Microsoft.Intune_RemoteAssistanceApp_Elevation",
	"Microsoft.Intune_MobileApps_ViewReports",
	"Microsoft.Intune_RemoteAssistance_Read",
	"Microsoft.Intune_EnrollmentProgramToken_Delete",
	"Microsoft.Intune_PartnerDeviceManagement_Modify",
	"Microsoft.Intune_DeviceCompliancePolices_Create",
	"Microsoft.Intune_DeviceEnrollmentManagers_Read",
	"Microsoft.Intune_SecurityBaselines_Delete",
	"Microsoft.Intune_RemoteTasks_LocateDevice",
	"Microsoft.Intune_RemoteTasks_ConfigurationManagerAction",
	"Microsoft.Intune_Organization_Read",
	"Microsoft.Intune_RemoteTasks_ResetPasscode",
	"Microsoft.Intune_SecurityBaselines_Read",
	"Microsoft.Intune_Roles_Delete",
	"Microsoft.Intune_EndpointAnalytics_Delete",
	"Microsoft.Intune_TelecomExpenses_Update",
	"Microsoft.Intune_Roles_Update",
	"Microsoft.Intune_RemoteAssistance_Update",
	"Microsoft.Intune_EndpointAnalytics_Update",
	"Microsoft.Intune_EndpointProtection_Read",
	"Microsoft.Intune_AssignmentFilter_Delete",
	"Microsoft.Intune_AssignmentFilter_Update",
	"Microsoft.Intune_MobileThreatDefense_Read",
	"Microsoft.Intune_RemoteAssistanceApp_TakeFullControl",
	"Microsoft.Intune_DeviceCompliancePolices_ViewReports",
	"Microsoft.Intune_PolicySets_Assign",
	"Microsoft.Intune_ManagedApps_Delete",
	"Microsoft.Intune_PolicySets_Create",
	"Microsoft.Intune_PolicySets_Update",
	"Microsoft.Intune_DeviceConfigurations_Read",
	"Microsoft.Intune_TermsAndConditions_Create",
	"Microsoft.Intune_AppleDeviceSerialNumbers_Delete",
	"Microsoft.Intune_AppleDeviceSerialNumbers_Update",
	"Microsoft.Intune_ManagedDevices_Read",
	"Microsoft.Intune_AppleDeviceSerialNumbers_Create",
	"Microsoft.Intune_Organization_Create",
	"Microsoft.Intune_Customization_Create",
	"Microsoft.Intune_RemoteTasks_GetFileVaultKey",
	"Microsoft.Intune_RemoteTasks_WindowsDefender",
	"Microsoft.Intune_RemoteTasks_RemoveDFCIManagement",
	"Microsoft.Intune_CloudAttach_ResourceExplorer",
	"Microsoft.Intune_CloudAttach_Timeline",
	"Microsoft.Intune_CloudAttach_Collections",
	"Microsoft.Intune_ManagedDeviceCleanupSettings_Update",
	"Microsoft.Intune_MultiAdminApproval_DeleteAccessPolicy",
	"Microsoft.Intune_MultiAdminApproval_UpdateAccessPolicy",
	"Microsoft.Intune_ServiceNow_UpdateConnector",
	"Microsoft.Intune_ServiceNow_ViewIncidents",
	"Microsoft.Intune_RemoteAssistanceApp_Unattended",
	"Microsoft.Intune_RemoteTasks_InitiateDeviceAttestation",
	"Microsoft.Intune_ManagedDevices_ReadBiosPassword",
	"Microsoft.Intune_EpmPolicy_ViewElevationRequests",
	"Microsoft.Intune_EpmPolicy_ModifyElevationRequests",
	"Microsoft.Intune_ManagedDeviceCleanupRule_Update",
	"Microsoft.Intune_ASRPolicy_ViewReports",
	"Microsoft.Intune_EDRPolicy_ViewReports",
	"Microsoft.Intune_EDRPolicy_Delete",
	"Microsoft.Intune_EDRPolicy_Assign",
	"Microsoft.Intune_ASRPolicy_Update",
	"Microsoft.Intune_ASRPolicy_Assign",
	"Microsoft.Intune_ASRPolicy_Delete",
	"Microsoft.Intune_EDRPolicy_Read",
	"Microsoft.Intune_AppControlPolicy_Assign",
	"Microsoft.Intune_AppleEnrollmentProfiles_RotateAppleAdminAccountPassword",
	"Microsoft.Intune_TermsAndConditions_Delete",
	"Microsoft.Intune_MobileApps_Read",
	"Microsoft.Intune_AndroidSync_Read",
	"Microsoft.Intune_ManagedApps_Read",
	"Microsoft.Intune_ChromebookSync_UpdateOnboarding",
	"Microsoft.Intune_ChromebookSync_DeleteOnboarding",
	"Microsoft.Intune_Organization_Update",
	"Microsoft.Intune_MicrosoftStoreForBusiness_Modify",
	"Microsoft.Intune_Organization_Delete",
	"Microsoft.Intune_ChromebookSync_Read",
	"Microsoft.Intune_RemoteTasks_CleanPC",
	"Microsoft.Intune_RemoteTasks_ShutDown",
	"Microsoft.Intune_SecurityTasks_Update",
	"Microsoft.Intune_TelecomExpenses_Read",
	"Microsoft.Intune_AndroidFota_Delete",
	"Microsoft.Intune_AndroidFota_Assign",
	"Microsoft.Intune_QuietTimePolicies_Create",
	"Microsoft.Intune_CloudAttach_Scripts",
	"Microsoft.Intune_CloudAttach_SoftwareUpdates",
	"Microsoft.Intune_OrganizationalMessages_Read",
	"Microsoft.Intune_OrganizationalMessages_Create",
	"Microsoft.Intune_RemoteTasks_OnDemandProactiveRemediation",
	"Microsoft.Intune_MultiAdminApproval_ReadAccessPolicy",
	"Microsoft.Intune_EpmPolicy_ViewReports",
	"Microsoft.Intune_OrganizationalMessages_UpdateControl",
	"Microsoft.Intune_RemoteTasks_PauseConfigurationRefresh",
	"Microsoft.Intune_CloudPki_Read",
	"Microsoft.Intune_TermsAndConditions_Read",
	"Microsoft.Intune_SecurityBaselines_Create",
	"Microsoft.Intune_ManagedGooglePlay_Modify",
	"Microsoft.Intune_SecurityBaselines_Update",
	"Microsoft.Intune_SecurityBaselines_Assign",
	"Microsoft.Intune_RemoteTasks_ManageSharedDeviceUsers",
	"Microsoft.Intune_RemoteTasks_RequestRemoteAssistance",
	"Microsoft.Intune_MobileApps_Delete",
	"Microsoft.Intune_ManagedApps_Assign",
	"Microsoft.Intune_MobileApps_Create",
	"Microsoft.Intune_RemoteTasks_Retire",
	"Microsoft.Intune_TermsAndConditions_Update",
	"Microsoft.Intune_DerivedCredentials_Modify",
	"Microsoft.Intune_MicrosoftDefenderATP_Read",
	"Microsoft.Intune_CertificateConnector_Read",
	"Microsoft.Intune_RemoteTasks_EnableLostMode",
	"Microsoft.Intune_PolicySets_Read",
	"Microsoft.Intune_ManagedApps_Wipe",
	"Microsoft.Intune_DeviceCompliancePolices_Read",
	"Microsoft.Intune_DeviceConfigurations_Assign",
	"Microsoft.Intune_AppleEnrollmentProfiles_Read",
	"Microsoft.Intune_Customization_Update",
	"Microsoft.Intune_PartnerDeviceManagement_Read",
	"Microsoft.Intune_DeviceConfigurations_Create",
	"Microsoft.Intune_EnrollmentProgramToken_Read",
	"Microsoft.Intune_AndroidSync_UpdateOnboarding",
	"Microsoft.Intune_RemoteTasks_DisableLostMode",
	"Microsoft.Intune_RemoteTasks_RotateBitLockerKeys",
	"Microsoft.Intune_MobileApps_Relate",
	"Microsoft.Intune_AndroidFota_Create",
	"Microsoft.Intune_QuietTimePolicies_Update",
	"Microsoft.Intune_CloudAttach_ScriptActions",
	"Microsoft.Intune_RemoteTasks_InitiateMDMKeyRecovery",
	"Microsoft.Intune_RemoteTasks_RotateLocalAdminPassword",
	"Microsoft.Intune_EpmPolicy_Read",
	"Microsoft.Intune_TenantAttachRecommendations_Read",
	"Microsoft.Intune_AndroidSync_UpdateEnrollmentProfiles",
	"Microsoft.Intune_ManagedDevices_Query",
	"Microsoft.Intune_AppControlPolicy_Read",
	"Microsoft.Intune_EDRPolicy_Update",
	"Microsoft.Intune_AppControlPolicy_ViewReports",
	"Microsoft.Intune_AppControlPolicy_Create",
	"Microsoft.Intune_AndroidSync_AndroidEnrollmentTimeMembershipAssign",
	"Microsoft.Intune_MultiAdminApproval_ApprovalForMultiAdminApproval",
	"Microsoft.Intune_AppleEnrollmentProfiles_ViewAppleAdminAccountPassword",
	"Microsoft.Intune_WindowsEnterpriseCertificate_Modify",
	"Microsoft.Intune_MobileApps_Assign",
	"Microsoft.Intune_ManagedApps_Update",
	"Microsoft.Intune_PolicySets_Delete",
	"Microsoft.Intune_ManagedApps_Create",
	"Microsoft.Intune_SecurityTasks_Read",
	"Microsoft.Intune_MobileThreatDefense_Modify",
	"Microsoft.Intune_RemoteTasks_Wipe",
	"Microsoft.Intune_DeviceConfigurations_ViewReports",
	"Microsoft.Intune_CorporateDeviceIdentifiers_Read",
	"Microsoft.Intune_Customization_Assign",
	"Microsoft.Intune_DeviceEnrollmentManagers_Update",
	"Microsoft.Intune_MicrosoftTunnelGateway_Read",
	"Microsoft.Intune_RemoteTasks_BypassActivationLock",
	"Microsoft.Intune_DeviceConfigurations_Update",
	"Microsoft.Intune_RemoteTasks_UpdateDeviceAccount",
	"Microsoft.Intune_AndroidFota_Read",
	"Microsoft.Intune_QuietTimePolicies_Delete",
	"Microsoft.Intune_QuietTimePolicies_Assign",
	"Microsoft.Intune_QuietTimePolicies_Read",
	"Microsoft.Intune_CloudAttach_CMPivot",
	"Microsoft.Intune_EpmPolicy_Create",
	"Microsoft.Intune_WindowsEnterpriseCertificate_Read",
	"Microsoft.Intune_Roles_Assign",
	"Microsoft.Intune_Reports_Read",
	"Microsoft.Intune_MobileApps_Update",
	"Microsoft.Intune_Customization_Read",
	"Microsoft.Intune_ManagedDevices_ViewReports",
	"Microsoft.Intune_RemoteTasks_SetDeviceName",
	"Microsoft.Intune_TermsAndConditions_Assign",
	"Microsoft.Intune_CertificateConnector_Modify",
	"Microsoft.Intune_DeviceConfigurations_Delete",
	"Microsoft.Intune_RemoteAssistance_ViewReports",
	"Microsoft.Intune_Customization_Delete",
	"Microsoft.Intune_AndroidFota_Update",
	"Microsoft.Intune_QuietTimePolicies_ViewReports",
	"Microsoft.Intune_CloudAttach_ApplicationActions",
	"Microsoft.Intune_CloudAttach_ClientDetails",
	"Microsoft.Intune_CloudAttach_Applications",
	"Microsoft.Intune_OrganizationalMessages_Delete",
	"Microsoft.Intune_OrganizationalMessages_Update",
	"Microsoft.Intune_CloudAttach_EnrollNow",
	"Microsoft.Intune_MultiAdminApproval_CreateAccessPolicy",
	"Microsoft.Intune_EpmPolicy_Assign",
	"Microsoft.Intune_EpmPolicy_Update",
	"Microsoft.Intune_EpmPolicy_Delete",
	"Microsoft.Intune_OrganizationalMessages_Assign",
	"Microsoft.Intune_WindowsOSRecovery_Assign",
	"Microsoft.Intune_EDRPolicy_Create",
	"Microsoft.Intune_AppControlPolicy_Delete",
	"Microsoft.Intune_DeviceConfigurations_UpdateWindowsRestore",
	"Microsoft.Intune_WindowsOSRecovery_Delete",
	"Microsoft.Intune_CloudPki_Update",
	"Microsoft.Intune_EnrollmentProfiles_EnrollmentTimeMembershipAssign",
	"Microsoft.Intune_ASRPolicy_Create",
	"Microsoft.Intune_ASRPolicy_Read",
	"Microsoft.Intune_AppControlPolicy_Update",
	"Microsoft.Intune_WindowsOSRecovery_Update",
	"Microsoft.Intune_WindowsOSRecovery_Create",
	"Microsoft.Intune_WindowsOSRecovery_Read",
	"Microsoft.Intune_CloudPki_Revoke",
	"Microsoft.Intune_CloudPki_Create",
	"Microsoft.Intune_RemoteTasks_ChangeAssignments",
	"Microsoft.Intune_RemoteTasks_Offboard",
	"Microsoft.Intune_EnrollmentProgramToken_ReleaseADEDevices",
}
resourceAction.SetAllowedResourceActions(allowedResourceActions)

resourceActions := []graphmodels.ResourceActionable {
	resourceAction,
}
rolePermission.SetResourceActions(resourceActions)

rolePermissions := []graphmodels.RolePermissionable {
	rolePermission,
}
requestBody.SetRolePermissions(rolePermissions)
roleScopeTagIds := []string {
	"0",
}
requestBody.SetRoleScopeTagIds(roleScopeTagIds)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
roleDefinitions, err := graphClient.DeviceManagement().RoleDefinitions().Post(context.Background(), requestBody, nil)


w365

POST https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions

// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewUnifiedRoleDefinition()
id := ""
requestBody.SetId(&id) 
description := "test"
requestBody.SetDescription(&description) 
displayName := "test"
requestBody.SetDisplayName(&displayName) 


unifiedRolePermission := graphmodels.NewUnifiedRolePermission()
condition := null
unifiedRolePermission.SetCondition(&condition) 
allowedResourceActions := []string {
	"Microsoft.CloudPC/CloudPCs/Read",
	"Microsoft.CloudPC/CloudPCs/Reprovision",
	"Microsoft.CloudPC/CloudPCs/Resize",
	"Microsoft.CloudPC/CloudPCs/EndGracePeriod",
	"Microsoft.CloudPC/CloudPCs/Restore",
	"Microsoft.CloudPC/CloudPCs/Reboot",
	"Microsoft.CloudPC/CloudPCs/Rename",
	"Microsoft.CloudPC/CloudPCs/Troubleshoot",
	"Microsoft.CloudPC/CloudPCs/ModifyDiskEncryptionType",
	"Microsoft.CloudPC/CloudPCs/ChangeUserAccountType",
	"Microsoft.CloudPC/CloudPCs/PlaceUnderReview",
	"Microsoft.CloudPC/CloudPCs/RetryPartnerAgentInstallation",
	"Microsoft.CloudPC/CloudPCs/ApplyCurrentProvisioningPolicy",
	"Microsoft.CloudPC/CloudPCs/CreateSnapshot",
	"Microsoft.CloudPC/CloudPCs/PowerOn",
	"Microsoft.CloudPC/CloudPCs/PowerOff",
	"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailover",
	"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailback",
	"Microsoft.CloudPC/CloudPCs/Start",
	"Microsoft.CloudPC/CloudPCs/Stop",
	"Microsoft.CloudPC/CloudPCs/GetCloudPcLaunchInfo",
	"Microsoft.CloudPC/CloudPCs/ReinstallAgent",
	"Microsoft.CloudPC/CloudPCs/CheckAgentStatus",
	"Microsoft.CloudPC/CloudPCs/RetrieveAgentStatus",
	"Microsoft.CloudPC/CloudPCs/Provision",
	"Microsoft.CloudPC/CloudPCs/Deprovision",
	"Microsoft.CloudPC/DeviceImages/Create",
	"Microsoft.CloudPC/DeviceImages/Delete",
	"Microsoft.CloudPC/DeviceImages/Read",
	"Microsoft.CloudPC/OnPremisesConnections/Create",
	"Microsoft.CloudPC/OnPremisesConnections/Delete",
	"Microsoft.CloudPC/OnPremisesConnections/Read",
	"Microsoft.CloudPC/OnPremisesConnections/Update",
	"Microsoft.CloudPC/OnPremisesConnections/RunHealthChecks",
	"Microsoft.CloudPC/OnPremisesConnections/UpdateAdDomainPassword",
	"Microsoft.CloudPC/ProvisioningPolicies/Assign",
	"Microsoft.CloudPC/ProvisioningPolicies/Apply",
	"Microsoft.CloudPC/ProvisioningPolicies/Create",
	"Microsoft.CloudPC/ProvisioningPolicies/Delete",
	"Microsoft.CloudPC/ProvisioningPolicies/Read",
	"Microsoft.CloudPC/ProvisioningPolicies/Update",
	"Microsoft.CloudPC/UserSettings/Assign",
	"Microsoft.CloudPC/UserSettings/Create",
	"Microsoft.CloudPC/UserSettings/Delete",
	"Microsoft.CloudPC/UserSettings/Read",
	"Microsoft.CloudPC/UserSettings/Update",
	"Microsoft.CloudPC/Roles/Read",
	"Microsoft.CloudPC/Roles/Create",
	"Microsoft.CloudPC/Roles/Update",
	"Microsoft.CloudPC/Roles/Delete",
	"Microsoft.CloudPC/RoleAssignments/Create",
	"Microsoft.CloudPC/RoleAssignments/Update",
	"Microsoft.CloudPC/RoleAssignments/Delete",
	"Microsoft.CloudPC/AuditData/Read",
	"Microsoft.CloudPC/SupportedRegion/Read",
	"Microsoft.CloudPC/ServicePlan/Read",
	"Microsoft.CloudPC/Snapshot/Read",
	"Microsoft.CloudPC/Snapshot/Share",
	"Microsoft.CloudPC/Snapshot/Import",
	"Microsoft.CloudPC/Snapshot/PurgeImportedSnapshot",
	"Microsoft.CloudPC/OrganizationSettings/Read",
	"Microsoft.CloudPC/OrganizationSettings/Update",
	"Microsoft.CloudPC/ExternalPartnerSettings/Read",
	"Microsoft.CloudPC/ExternalPartnerSettings/Create",
	"Microsoft.CloudPC/ExternalPartnerSettings/Update",
	"Microsoft.CloudPC/PerformanceReports/Read",
	"Microsoft.CloudPC/SharedUseServicePlans/Read",
	"Microsoft.CloudPC/FrontLineServicePlans/Read",
	"Microsoft.CloudPC/SharedUseLicenseUsageReports/Read",
	"Microsoft.CloudPC/FrontlineReports/Read",
	"Microsoft.CloudPC/CrossRegionDisasterRecovery/Read",
	"Microsoft.CloudPC/BulkActions/Read",
	"Microsoft.CloudPC/BulkActions/Write",
	"Microsoft.CloudPC/ActionStatus/Read",
	"Microsoft.CloudPC/InaccessibleReports/Read",
	"Microsoft.CloudPC/MaintenanceWindows/Assign",
	"Microsoft.CloudPC/MaintenanceWindows/Create",
	"Microsoft.CloudPC/MaintenanceWindows/Delete",
	"Microsoft.CloudPC/MaintenanceWindows/Read",
	"Microsoft.CloudPC/MaintenanceWindows/Update",
	"Microsoft.CloudPC/DeviceRecommendation/Read",
	"Microsoft.CloudPC/CloudApps/Read",
	"Microsoft.CloudPC/CloudApps/Publish",
	"Microsoft.CloudPC/CloudApps/Update",
	"Microsoft.CloudPC/CloudApps/Reset",
	"Microsoft.CloudPC/CloudApps/Unpublish",
	"Microsoft.CloudPC/Settings/Assign",
	"Microsoft.CloudPC/Settings/Create",
	"Microsoft.CloudPC/Settings/Read",
	"Microsoft.CloudPC/Settings/Update",
	"Microsoft.CloudPC/Settings/Delete",
	"Microsoft.CloudPC/AdminHighlights/Operate",
}
unifiedRolePermission.SetAllowedResourceActions(allowedResourceActions)

rolePermissions := []graphmodels.UnifiedRolePermissionable {
	unifiedRolePermission,
}
requestBody.SetRolePermissions(rolePermissions)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
roleDefinitions, err := graphClient.RoleManagement().CloudPC().RoleDefinitions().Post(context.Background(), requestBody, nil)


windows autopatch

{"Name":"test",
"Description":"test",
"Permissions":[
	"Microsoft.Autopatch_Roles_Read","Microsoft.Autopatch_Roles_Create","Microsoft.Autopatch_Roles_Edit","Microsoft.Autopatch_Roles_Delete","Microsoft.Autopatch_Roles_Assign","Microsoft.Autopatch_Reports_Read","Microsoft.Autopatch_Reports_ActionDiscoverDevices","Microsoft.Autopatch_Reports_ActionAssignRing","Microsoft.Autopatch_Reports_ActionExcludeDevices","Microsoft.Autopatch_Reports_ActionRestoreExcludedDevices","Microsoft.Autopatch_AutopatchGroups_Read","Microsoft.Autopatch_AutopatchGroups_Create","Microsoft.Autopatch_AutopatchGroups_Edit","Microsoft.Autopatch_AutopatchGroups_Delete","Microsoft.Autopatch_Messages_Read","Microsoft.Autopatch_SupportRequests_Read"
],
"ScopeTags":["0"]
}