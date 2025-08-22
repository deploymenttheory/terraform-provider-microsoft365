# Microsoft 365 Device Management Resources - Development Status

This document tracks the implementation status of Microsoft Graph Beta Device Management API endpoints as Terraform resources.

| API Endpoint | Provider Resource Name | Implementation Status | Resource Supported Assignment Types |
|--------------|------------------------|----------------------|-----------------------------------|
| `/deviceManagement/advancedThreatProtectionOnboardingStateSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidDeviceOwnerEnrollmentProfiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidForWorkAppConfigurationSchemas` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidForWorkEnrollmentProfiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidForWorkSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidManagedStoreAccountEnterpriseSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/androidManagedStoreAppConfigurationSchemas` | - | ❌ Not Implemented | - |
| `/deviceManagement/applePushNotificationCertificate` | - | ❌ Not Implemented | - |
| `/deviceManagement/appleUserInitiatedEnrollmentProfiles` | `microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment` | ✅ Implemented | TBD |
| `/deviceManagement/assignmentFilters` | `microsoft365_graph_beta_device_management_assignment_filter` | ✅ Implemented | - |
| `/deviceManagement/auditEvents` | - | ❌ Not Implemented | - |
| `/deviceManagement/autopilotEvents` | - | ❌ Not Implemented | - |
| `/deviceManagement/cartToClassAssociations` | - | ❌ Not Implemented | - |
| `/deviceManagement/categories` | `microsoft365_graph_beta_device_management_device_category` | ✅ Implemented | - |
| `/deviceManagement/certificateConnectorDetails` | - | ❌ Not Implemented | - |
| `/deviceManagement/chromeOSOnboardingSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/cloudCertificationAuthority` | - | ❌ Not Implemented | - |
| `/deviceManagement/cloudCertificationAuthorityLeafCertificate` | - | ❌ Not Implemented | - |
| `/deviceManagement/cloudPCConnectivityIssues` | - | ❌ Not Implemented | - |
| `/deviceManagement/comanagedDevices` | - | ❌ Not Implemented | - |
| `/deviceManagement/comanagementEligibleDevices` | - | ❌ Not Implemented | - |
| `/deviceManagement/complianceCategories` | - | ❌ Not Implemented | - |
| `/deviceManagement/complianceManagementPartners` | - | ❌ Not Implemented | - |
| `/deviceManagement/compliancePolicies` | Multiple compliance policy resources | ✅ Implemented | All Groups + Filters |
| `/deviceManagement/complianceSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/conditionalAccessSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/configManagerCollections` | - | ❌ Not Implemented | - |
| `/deviceManagement/configurationCategories` | - | ❌ Not Implemented | - |
| `/deviceManagement/configurationPolicies` | `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | ✅ Implemented | All Groups + Filters |
| `/deviceManagement/configurationPolicyTemplates` | `microsoft365_graph_beta_device_management_settings_catalog_template_json` | ✅ Implemented | All Groups + Filters |
| `/deviceManagement/configurationSettings` | `microsoft365_graph_beta_device_management_reuseable_policy_settings` | ✅ Implemented | TBD |
| `/deviceManagement/dataSharingConsents` | - | ❌ Not Implemented | - |
| `/deviceManagement/depOnboardingSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/derivedCredentials` | - | ❌ Not Implemented | - |
| `/deviceManagement/detectedApps` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceCategories` | `microsoft365_graph_beta_device_management_device_category` | ✅ Implemented | - |
| `/deviceManagement/deviceCompliancePolicies` | Multiple compliance policy resources | ✅ Implemented | All Groups + Filters |
| `/deviceManagement/deviceCompliancePolicyDeviceStateSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceCompliancePolicySettingStateSummaries` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceComplianceScripts` | `microsoft365_graph_beta_device_management_windows_device_compliance_script` | ✅ Implemented | TBD |
| `/deviceManagement/deviceConfigurationConflictSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceConfigurationDeviceStateSummaries` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceConfigurationRestrictedAppsViolations` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceConfigurations` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceConfigurationsAllManagedDeviceCertificateStates` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceConfigurationUserStateSummaries` | - | ❌ Not Implemented | - |
| `/deviceManagement/deviceCustomAttributeShellScripts` | `microsoft365_graph_beta_device_management_macos_custom_attribute_script` | ✅ Implemented | TBD |
| `/deviceManagement/deviceEnrollmentConfigurations` | `microsoft365_graph_beta_device_management_windows_enrollment_status_page` | ✅ Implemented | All Licensed Users + All Devices + Inclusion Groups |
| `/deviceManagement/deviceHealthScripts` | `microsoft365_graph_beta_device_management_windows_remediation_script` | ✅ Implemented | TBD |
| `/deviceManagement/deviceManagementScripts` | `microsoft365_graph_beta_device_management_windows_platform_script` | ✅ Implemented | TBD |
| `/deviceManagement/deviceShellScripts` | `microsoft365_graph_beta_device_management_macos_platform_script` | ✅ Implemented | TBD |
| `/deviceManagement/domainJoinConnectors` | - | ❌ Not Implemented | - |
| `/deviceManagement/elevationRequests` | - | ❌ Not Implemented | - |
| `/deviceManagement/embeddedSIMActivationCodePools` | - | ❌ Not Implemented | - |
| `/deviceManagement/enableAndroidDeviceAdministratorEnrollment` | - | ❌ Not Implemented | - |
| `/deviceManagement/enableEndpointPrivilegeManagement` | `microsoft365_graph_beta_device_management_endpoint_privilege_management_json` | ✅ Implemented | All Groups + Filters |
| `/deviceManagement/enableLegacyPcManagement` | - | ❌ Not Implemented | - |
| `/deviceManagement/enableUnlicensedAdminstrators` | - | ❌ Not Implemented | - |
| `/deviceManagement/endpointPrivilegeManagementProvisioningStatus` | - | ❌ Not Implemented | - |
| `/deviceManagement/evaluateAssignmentFilter` | - | ❌ Not Implemented | - |
| `/deviceManagement/exchangeConnectors` | - | ❌ Not Implemented | - |
| `/deviceManagement/exchangeOnPremisesPolicies` | - | ❌ Not Implemented | - |
| `/deviceManagement/exchangeOnPremisesPolicy` | - | ❌ Not Implemented | - |
| `/deviceManagement/getAssignedRoleDetails` | - | ❌ Not Implemented | - |
| `/deviceManagement/getAssignmentFiltersStatusDetails` | - | ❌ Not Implemented | - |
| `/deviceManagement/getComanagedDevicesSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/getComanagementEligibleDevicesSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/getEffectivePermissions` | - | ❌ Not Implemented | - |
| `/deviceManagement/getEffectivePermissionsWithScope` | - | ❌ Not Implemented | - |
| `/deviceManagement/getRoleScopeTagsByIdsWithIds` | - | ❌ Not Implemented | - |
| `/deviceManagement/getRoleScopeTagsByResourceWithResource` | - | ❌ Not Implemented | - |
| `/deviceManagement/getSuggestedEnrollmentLimitWithEnrollmentType` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyCategories` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyConfigurations` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyDefinitionFiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyDefinitions` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyMigrationReports` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyObjectFiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/groupPolicyUploadedDefinitionFiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/hardwareConfigurations` | - | ❌ Not Implemented | - |
| `/deviceManagement/hardwarePasswordDetails` | - | ❌ Not Implemented | - |
| `/deviceManagement/hardwarePasswordInfo` | - | ❌ Not Implemented | - |
| `/deviceManagement/importedDeviceIdentities` | - | ❌ Not Implemented | - |
| `/deviceManagement/importedWindowsAutopilotDeviceIdentities` | - | ❌ Not Implemented | - |
| `/deviceManagement/intents` | - | ❌ Not Implemented | - |
| `/deviceManagement/intuneBrandingProfiles` | `microsoft365_graph_beta_device_management_intune_branding_profile` | 🚧 Work In Progress | TBD |
| `/deviceManagement/iosUpdateStatuses` | - | ❌ Not Implemented | - |
| `/deviceManagement/macOSSoftwareUpdateAccountSummaries` | `microsoft365_graph_beta_device_management_macos_software_update_configuration` | ✅ Implemented | All Groups |
| `/deviceManagement/managedDeviceCleanupRules` | `microsoft365_graph_beta_device_management_managed_device_cleanup_rule` | 🚧 Work In Progress | - |
| `/deviceManagement/managedDeviceEncryptionStates` | - | ❌ Not Implemented | - |
| `/deviceManagement/managedDeviceOverview` | - | ❌ Not Implemented | - |
| `/deviceManagement/managedDevices` | - | ❌ Not Implemented | - |
| `/deviceManagement/managedDeviceWindowsOSImages` | - | ❌ Not Implemented | - |
| `/deviceManagement/microsoftTunnelConfigurations` | - | ❌ Not Implemented | - |
| `/deviceManagement/microsoftTunnelHealthThresholds` | - | ❌ Not Implemented | - |
| `/deviceManagement/microsoftTunnelServerLogCollectionResponses` | - | ❌ Not Implemented | - |
| `/deviceManagement/microsoftTunnelSites` | - | ❌ Not Implemented | - |
| `/deviceManagement/mobileAppTroubleshootingEvents` | - | ❌ Not Implemented | - |
| `/deviceManagement/mobileThreatDefenseConnectors` | - | ❌ Not Implemented | - |
| `/deviceManagement/monitoring` | - | ❌ Not Implemented | - |
| `/deviceManagement/ndesConnectors` | - | ❌ Not Implemented | - |
| `/deviceManagement/notificationMessageTemplates` | `microsoft365_graph_beta_device_management_device_compliance_notification_template` | ✅ Implemented | - |
| `/deviceManagement/operationApprovalPolicies` | `microsoft365_graph_beta_device_management_operation_approval_policy` | ✅ Implemented | TBD |
| `/deviceManagement/operationApprovalRequests` | - | ❌ Not Implemented | - |
| `/deviceManagement/privilegeManagementElevations` | - | ❌ Not Implemented | - |
| `/deviceManagement/remoteActionAudits` | - | ❌ Not Implemented | - |
| `/deviceManagement/remoteAssistancePartners` | - | ❌ Not Implemented | - |
| `/deviceManagement/remoteAssistanceSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/reports` | - | ❌ Not Implemented | - |
| `/deviceManagement/resourceAccessProfiles` | - | ❌ Not Implemented | - |
| `/deviceManagement/resourceOperations` | `microsoft365_graph_beta_device_management_rbac_resource_operation` | ✅ Implemented | - |
| `/deviceManagement/retrieveUserRoleDetailWithUserid` | - | ❌ Not Implemented | - |
| `/deviceManagement/reusablePolicySettings` | `microsoft365_graph_beta_device_management_reuseable_policy_settings` | ✅ Implemented | TBD |
| `/deviceManagement/reusableSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/roleAssignments` | `microsoft365_graph_beta_device_management_role_assignment` | ✅ Implemented | TBD |
| `/deviceManagement/roleDefinitions` | `microsoft365_graph_beta_device_management_role_definition` | ✅ Implemented | - |
| `/deviceManagement/roleScopeTags` | `microsoft365_graph_beta_device_management_role_scope_tag` | ✅ Implemented | Groups Only |
| `/deviceManagement/scopedForResourceWithResource` | - | ❌ Not Implemented | - |
| `/deviceManagement/sendCustomNotificationToCompanyPortal` | - | ❌ Not Implemented | - |
| `/deviceManagement/serviceNowConnections` | - | ❌ Not Implemented | - |
| `/deviceManagement/settingDefinitions` | - | ❌ Not Implemented | - |
| `/deviceManagement/softwareUpdateStatusSummary` | - | ❌ Not Implemented | - |
| `/deviceManagement/telecomExpenseManagementPartners` | - | ❌ Not Implemented | - |
| `/deviceManagement/templateInsights` | - | ❌ Not Implemented | - |
| `/deviceManagement/templates` | - | ❌ Not Implemented | - |
| `/deviceManagement/templateSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/tenantAttachRBAC` | - | ❌ Not Implemented | - |
| `/deviceManagement/termsAndConditions` | `microsoft365_graph_beta_device_management_terms_and_conditions` | ✅ Implemented | All Licensed Users + Groups + SCCM Collections |
| `/deviceManagement/troubleshootingEvents` | - | ❌ Not Implemented | - |
| `/deviceManagement/windowsAutopilotDeploymentProfiles` | `microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile` | ✅ Implemented | TBD |
| `/deviceManagement/windowsAutopilotDeviceIdentities` | `microsoft365_graph_beta_device_management_windows_autopilot_device_identity` | ✅ Implemented | - |
| `/deviceManagement/windowsAutopilotSettings` | - | ❌ Not Implemented | - |
| `/deviceManagement/windowsDriverUpdateProfiles` | `microsoft365_graph_beta_device_management_windows_driver_update_profile` | ✅ Implemented | TBD |
| `/deviceManagement/windowsFeatureUpdateProfiles` | `microsoft365_graph_beta_device_management_windows_feature_update_profile` | ✅ Implemented | TBD |
| `/deviceManagement/windowsInformationProtectionAppLearningSummaries` | - | ❌ Not Implemented | - |
| `/deviceManagement/windowsInformationProtectionNetworkLearningSummaries` | - | ❌ Not Implemented | - |
| `/deviceManagement/windowsMalwareInformation` | - | ❌ Not Implemented | - |
| `/deviceManagement/windowsQualityUpdatePolicies` | `microsoft365_graph_beta_device_management_windows_quality_update_policy` | ✅ Implemented | TBD |
| `/deviceManagement/windowsQualityUpdateProfiles` | `microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy` | ✅ Implemented | TBD |
| `/deviceManagement/windowsUpdateCatalogItems` | `microsoft365_graph_beta_device_management_windows_driver_update_inventory` | ✅ Implemented | - |
| `/deviceManagement/zebraFotaArtifacts` | - | ❌ Not Implemented | - |
| `/deviceManagement/zebraFotaConnector` | - | ❌ Not Implemented | - |
| `/deviceManagement/zebraFotaDeployments` | - | ❌ Not Implemented | - |

## Assignment Types Legend

- **All Groups + Filters**: Supports all assignment types with group filters (allDevicesAssignmentTarget, allLicensedUsersAssignmentTarget, groupAssignmentTarget, exclusionGroupAssignmentTarget)
- **All Groups**: Supports all assignment types without filters (allDevicesAssignmentTarget, allLicensedUsersAssignmentTarget, groupAssignmentTarget, exclusionGroupAssignmentTarget)
- **All Licensed Users + All Devices + Groups**: Supports allLicensedUsersAssignmentTarget, allDevicesAssignmentTarget, and groupAssignmentTarget (no exclusions)
- **All Licensed Users + Groups + SCCM Collections**: Supports allLicensedUsersAssignmentTarget, groupAssignmentTarget, and configurationManagerCollection
- **Groups Only**: Supports groupAssignmentTarget only
- **TBD**: Assignment types need to be determined/documented
- **-**: No assignments supported (read-only or configuration resources)

## Status Legend

- ✅ **Implemented**: Resource is fully implemented and available
- 🚧 **Work In Progress**: Resource is partially implemented or has known issues
- ❌ **Not Implemented**: Resource has not been implemented yet

## Implementation Statistics

- **Total Endpoints**: 135
- **Implemented**: 33 (24.4%)
- **Work In Progress**: 2 (1.5%)
- **Not Implemented**: 100 (74.1%)

---
*Last updated: August 2025*