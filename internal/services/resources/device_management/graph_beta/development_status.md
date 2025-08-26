# Microsoft 365 Device Management Resources - Development Status

This document tracks the implementation status of Microsoft Graph Beta Device Management API endpoints as Terraform resources.

| Provider Resource Name | Release Version | API Endpoint | Test Harness | Implementation Status | Resource Supported Assignment Types |
|------------------------|----------------|--------------|--------------|----------------------|-----------------------------------|
| - | - | `/deviceManagement/advancedThreatProtectionOnboardingStateSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidDeviceOwnerEnrollmentProfiles` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidForWorkAppConfigurationSchemas` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidForWorkEnrollmentProfiles` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidForWorkSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidManagedStoreAccountEnterpriseSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/androidManagedStoreAppConfigurationSchemas` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/applePushNotificationCertificate` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment` | v0.26.0-alpha | `/deviceManagement/appleUserInitiatedEnrollmentProfiles` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_assignment_filter` | v0.25.0-alpha | `/deviceManagement/assignmentFilters` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| - | - | `/deviceManagement/auditEvents` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/autopilotEvents` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/cartToClassAssociations` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_device_category` | v0.24.0-alpha | `/deviceManagement/categories` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| - | - | `/deviceManagement/certificateConnectorDetails` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/chromeOSOnboardingSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/cloudCertificationAuthority` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/cloudCertificationAuthorityLeafCertificate` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/cloudPCConnectivityIssues` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/comanagedDevices` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/comanagementEligibleDevices` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/complianceCategories` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/complianceManagementPartners` | - | ❌ Not Implemented | - |
| Multiple compliance policy resources | v0.24.0-alpha | `/deviceManagement/compliancePolicies` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| - | - | `/deviceManagement/complianceSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/conditionalAccessSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/configManagerCollections` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/configurationCategories` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | v0.24.0-alpha | `/deviceManagement/configurationPolicies` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| `microsoft365_graph_beta_device_management_settings_catalog_template_json` | v0.25.0-alpha | `/deviceManagement/configurationPolicyTemplates` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| `microsoft365_graph_beta_device_management_reuseable_policy_settings` | v0.25.0-alpha | `/deviceManagement/configurationSettings` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/dataSharingConsents` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/depOnboardingSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/derivedCredentials` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/detectedApps` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_device_category` | v0.24.0-alpha | `/deviceManagement/deviceCategories` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| Multiple compliance policy resources | v0.24.0-alpha | `/deviceManagement/deviceCompliancePolicies` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| - | - | `/deviceManagement/deviceCompliancePolicyDeviceStateSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceCompliancePolicySettingStateSummaries` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_windows_device_compliance_script` | v0.25.0-alpha | `/deviceManagement/deviceComplianceScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/deviceConfigurationConflictSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceConfigurationDeviceStateSummaries` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceConfigurationRestrictedAppsViolations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceConfigurations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceConfigurationsAllManagedDeviceCertificateStates` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/deviceConfigurationUserStateSummaries` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_macos_custom_attribute_script` | v0.25.0-alpha | `/deviceManagement/deviceCustomAttributeShellScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_enrollment_status_page` | v0.25.0-alpha | `/deviceManagement/deviceEnrollmentConfigurations` | ✅ | ✅ Implemented Unit and Acceptance tests | All Licensed Users + All Devices + Inclusion Groups |
| `microsoft365_graph_beta_device_management_windows_remediation_script` | v0.25.0-alpha | `/deviceManagement/deviceHealthScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_platform_script` | v0.25.0-alpha | `/deviceManagement/deviceManagementScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_macos_platform_script` | v0.25.0-alpha | `/deviceManagement/deviceShellScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/domainJoinConnectors` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/elevationRequests` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/embeddedSIMActivationCodePools` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/enableAndroidDeviceAdministratorEnrollment` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_endpoint_privilege_management_json` | v0.26.0-alpha | `/deviceManagement/enableEndpointPrivilegeManagement` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| - | - | `/deviceManagement/enableLegacyPcManagement` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/enableUnlicensedAdminstrators` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/endpointPrivilegeManagementProvisioningStatus` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/evaluateAssignmentFilter` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/exchangeConnectors` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/exchangeOnPremisesPolicies` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/exchangeOnPremisesPolicy` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getAssignedRoleDetails` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getAssignmentFiltersStatusDetails` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getComanagedDevicesSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getComanagementEligibleDevicesSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getEffectivePermissions` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getEffectivePermissionsWithScope` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getRoleScopeTagsByIdsWithIds` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getRoleScopeTagsByResourceWithResource` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/getSuggestedEnrollmentLimitWithEnrollmentType` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyCategories` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyConfigurations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyDefinitionFiles` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyDefinitions` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyMigrationReports` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyObjectFiles` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/groupPolicyUploadedDefinitionFiles` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/hardwareConfigurations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/hardwarePasswordDetails` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/hardwarePasswordInfo` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/importedDeviceIdentities` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/importedWindowsAutopilotDeviceIdentities` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/intents` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_intune_branding_profile` | v0.27.0-alpha | `/deviceManagement/intuneBrandingProfiles` | ✅ | 🚧 Work In Progress | TBD |
| - | - | `/deviceManagement/iosUpdateStatuses` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_macos_software_update_configuration` | v0.26.0-alpha | `/deviceManagement/macOSSoftwareUpdateAccountSummaries` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups |
| `microsoft365_graph_beta_device_management_managed_device_cleanup_rule` | v0.27.0-alpha | `/deviceManagement/managedDeviceCleanupRules` | ✅ | 🚧 Work In Progress | - |
| - | - | `/deviceManagement/managedDeviceEncryptionStates` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/managedDeviceOverview` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/managedDevices` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/managedDeviceWindowsOSImages` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/microsoftTunnelConfigurations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/microsoftTunnelHealthThresholds` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/microsoftTunnelServerLogCollectionResponses` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/microsoftTunnelSites` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/mobileAppTroubleshootingEvents` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/mobileThreatDefenseConnectors` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/monitoring` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/ndesConnectors` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_device_compliance_notification_template` | v0.26.0-alpha | `/deviceManagement/notificationMessageTemplates` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| `microsoft365_graph_beta_device_management_operation_approval_policy` | v0.26.0-alpha | `/deviceManagement/operationApprovalPolicies` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/operationApprovalRequests` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/privilegeManagementElevations` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/remoteActionAudits` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/remoteAssistancePartners` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/remoteAssistanceSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/reports` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/resourceAccessProfiles` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_rbac_resource_operation` | v0.25.0-alpha | `/deviceManagement/resourceOperations` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| - | - | `/deviceManagement/retrieveUserRoleDetailWithUserid` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_reuseable_policy_settings` | v0.25.0-alpha | `/deviceManagement/reusablePolicySettings` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/reusableSettings` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_role_assignment` | v0.25.0-alpha | `/deviceManagement/roleAssignments` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_role_definition` | v0.25.0-alpha | `/deviceManagement/roleDefinitions` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| `microsoft365_graph_beta_device_management_role_scope_tag` | v0.25.0-alpha | `/deviceManagement/roleScopeTags` | ✅ | ✅ Implemented Unit and Acceptance tests | Groups Only |
| - | - | `/deviceManagement/scopedForResourceWithResource` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/sendCustomNotificationToCompanyPortal` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/serviceNowConnections` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/settingDefinitions` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/softwareUpdateStatusSummary` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/telecomExpenseManagementPartners` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/templateInsights` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/templates` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/templateSettings` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/tenantAttachRBAC` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_terms_and_conditions` | v0.25.0-alpha | `/deviceManagement/termsAndConditions` | ✅ | ✅ Implemented Unit and Acceptance tests | All Licensed Users + Groups + SCCM Collections |
| - | - | `/deviceManagement/troubleshootingEvents` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile` | v0.26.0-alpha | `/deviceManagement/windowsAutopilotDeploymentProfiles` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_autopilot_device_identity` | v0.26.0-alpha | `/deviceManagement/windowsAutopilotDeviceIdentities` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| - | - | `/deviceManagement/windowsAutopilotSettings` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_windows_driver_update_profile` | v0.27.0-alpha | `/deviceManagement/windowsDriverUpdateProfiles` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_feature_update_profile` | v0.27.0-alpha | `/deviceManagement/windowsFeatureUpdateProfiles` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| - | - | `/deviceManagement/windowsInformationProtectionAppLearningSummaries` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/windowsInformationProtectionNetworkLearningSummaries` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/windowsMalwareInformation` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_windows_quality_update_policy` | v0.27.0-alpha | `/deviceManagement/windowsQualityUpdatePolicies` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy` | v0.27.0-alpha | `/deviceManagement/windowsQualityUpdateProfiles` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_windows_driver_update_inventory` | v0.26.0-alpha | `/deviceManagement/windowsUpdateCatalogItems` | ✅ | ✅ Implemented Unit and Acceptance tests | - |
| - | - | `/deviceManagement/zebraFotaArtifacts` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/zebraFotaConnector` | - | ❌ Not Implemented | - |
| - | - | `/deviceManagement/zebraFotaDeployments` | - | ❌ Not Implemented | - |
| `microsoft365_graph_beta_device_management_autopatch_groups` | v0.27.0-alpha | `/deviceManagement/autopatchGroups` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_linux_platform_script` | v0.27.0-alpha | ``/deviceManagement/configurationPolicies`` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_linux_device_compliance_script` | v0.27.0-alpha | `/deviceManagement/linuxDeviceComplianceScripts` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_device_enrollment_notification` | v0.27.0-alpha | `/deviceManagement/deviceEnrollmentNotifications` | ✅ | ✅ Implemented Unit and Acceptance tests | TBD |
| `microsoft365_graph_beta_device_management_macos_device_configuration_templates` | v0.27.0-alpha | `/deviceManagement/macosDeviceConfigurationTemplates` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| `microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls` | v0.27.0-alpha | `/deviceManagement/appControlForBusinessBuiltInControls` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |
| `microsoft365_graph_beta_device_management_app_control_for_business_managed_installer` | v0.27.0-alpha | `/deviceManagement/appControlForBusinessManagedInstaller` | ✅ | ✅ Implemented Unit and Acceptance tests | All Groups + Filters |

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

## Test Harness Legend
- ✅ **Tested**: Resource has automated tests (unit and acceptance tests)
- ❌ **Not Tested**: Resource does not have automated tests
- 🚧 **Partial Testing**: Resource has some tests but coverage is incomplete

## Implementation Statistics

- **Total Endpoints**: 142
- **Implemented**: 40 (28.2%)
- **Work In Progress**: 2 (1.4%)
- **Not Implemented**: 100 (70.4%)

---
*Last updated: August 2025*