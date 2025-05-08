---
page_title: "Provider Roadmap"
description: |-
  Development roadmap for the Microsoft 365 Terraform Provider.
---

# Microsoft 365 Provider Roadmap

This document outlines the development roadmap for the Microsoft 365 Terraform Provider. It provides visibility into our implementation plans and priorities for upcoming releases.

## Understanding the Roadmap

This roadmap reflects our current planning and is subject to change based on customer feedback, technical challenges, and changes to the Microsoft Graph API. 

### Column Definitions

- **Resource Name**: The name of the Microsoft Graph resource as it appears in the API documentation.
- **Category**: The logical grouping of the resource (Device Management, App Management, Identity & Access, etc.).
- **API Path**: Whether the resource is available in the Microsoft Graph v1.0 endpoint or beta endpoint.
- **Status**: Current implementation status:
  - **Planned**: On our roadmap but work has not yet begun
  - **In Progress**: Development is currently underway
  - **Completed**: Implementation is finished and available in a released version
  - **Investigating**: We're researching the feasibility and approach
  - **Backlog**: A future consideration not actively planned for near-term releases
- **Priority**: The relative importance of this resource:
  - **High**: Critical functionality targeted for the next release
  - **Medium**: Important functionality planned for near-term releases
  - **Low**: Desirable but not prioritized for upcoming releases
- **Target Version**: The planned provider version where this resource will be available.
- **Dependencies**: Other resources or features that must be implemented first.
- **Notes**: Additional context, limitations, or implementation details.

## Development Status

### Intune Device Configuration Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| DeviceConfigurations | Device Configuration | Beta | In Progress | High | v1.2.0 | None | Complex schema |
| DeviceConfigurationConflictSummary | Device Configuration | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceConfigurationDeviceStateSummaries | Device Configuration | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceConfigurationRestrictedAppsViolations | Device Configuration | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceConfigurationsAllManagedDeviceCertificateStates | Device Configuration | Beta | Planned | Medium | v1.3.0 | DeviceConfigurations | |
| DeviceConfigurationUserStateSummaries | Device Configuration | Beta | Planned | Medium | v1.3.0 | DeviceConfigurations | |
| ConfigurationCategories | Device Configuration | Beta | Planned | Medium | v1.3.0 | None | |
| ConfigurationPolicies | Device Configuration | Beta | Planned | High | v1.2.0 | None | |
| ConfigurationPolicyTemplates | Device Configuration | Beta | Planned | High | v1.2.0 | None | |
| ConfigurationSettings | Device Configuration | Beta | Planned | High | v1.2.0 | None | |
| AssignmentFilters | Device Configuration | Beta | Planned | High | v1.2.0 | None | Core functionality |

### Intune Device Compliance Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| DeviceCompliancePolicies | Device Compliance | Beta | Planned | High | v1.2.0 | None | |
| DeviceCompliancePolicyDeviceStateSummary | Device Compliance | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceCompliancePolicySettingStateSummaries | Device Compliance | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceComplianceScripts | Device Compliance | Beta | Planned | Medium | v1.3.0 | None | |
| ComplianceCategories | Device Compliance | Beta | Planned | Medium | v1.3.0 | None | |
| CompliancePolicies | Device Compliance | Beta | Planned | High | v1.2.0 | None | |
| ComplianceSettings | Device Compliance | Beta | Planned | High | v1.2.0 | None | |
| ComplianceManagementPartners | Device Compliance | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Device Enrollment Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| DeviceEnrollmentConfigurations | Device Enrollment | Beta | Planned | High | v1.2.0 | None | |
| AndroidDeviceOwnerEnrollmentProfiles | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidForWorkEnrollmentProfiles | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| AppleUserInitiatedEnrollmentProfiles | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| DepOnboardingSettings | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsAutopilotDeploymentProfiles | Device Enrollment | Beta | Planned | High | v1.2.0 | None | |
| WindowsAutopilotDeviceIdentities | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsAutopilotSettings | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |
| AutopilotEvents | Device Enrollment | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Windows Update Management Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| WindowsQualityUpdateProfiles | Windows Updates | Beta | Completed | Medium | v1.1.0 | None | |
| WindowsDriverUpdateProfiles | Windows Updates | Beta | Completed | Medium | v1.1.0 | None | |
| WindowsFeatureUpdateProfiles | Windows Updates | Beta | Completed | Medium | v1.1.0 | None | |
| WindowsUpdateCatalogItems | Windows Updates | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsQualityUpdatePolicies | Windows Updates | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Device Management Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| ManagedDevices | Device Management | Beta | Planned | High | v1.2.0 | None | |
| ComanagedDevices | Device Management | Beta | Planned | High | v1.2.0 | None | |
| ComanagementEligibleDevices | Device Management | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceCategories | Device Management | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceManagementPartners | Device Management | Beta | Planned | Medium | v1.3.0 | None | |
| DetectedApps | Device Management | Beta | Planned | Medium | v1.3.0 | None | |
| ManagedDeviceOverview | Device Management | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Device Scripts and Health Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| DeviceHealthScripts | Device Scripts | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceManagementScripts | Device Scripts | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceShellScripts | Device Scripts | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceCustomAttributeShellScripts | Device Scripts | Beta | Planned | Medium | v1.4.0 | None | |

### Intune Platform Integration Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| ApplePushNotificationCertificate | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidForWorkSettings | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidForWorkAppConfigurationSchemas | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidManagedStoreAccountEnterpriseSettings | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidManagedStoreAppConfigurationSchemas | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| ChromeOSOnboardingSettings | Platform Integration | Beta | Planned | Medium | v1.4.0 | None | |
| CertificateConnectorDetails | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| ExchangeConnectors | Platform Integration | Beta | Planned | Medium | v1.3.0 | None | |
| ConditionalAccessSettings | Platform Integration | Beta | Planned | High | v1.2.0 | None | |

### Intune Advanced Security Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| AdvancedThreatProtectionOnboardingStateSummary | Security | Beta | Planned | Low | v1.4.0 | None | |
| WindowsMalwareInformation | Security | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsInformationProtectionAppLearningSummaries | Security | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsInformationProtectionNetworkLearningSummaries | Security | Beta | Planned | Medium | v1.3.0 | None | |
| DerivedCredentials | Security | Beta | Planned | Medium | v1.4.0 | None | |

### Intune App Management Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| MobileApps | App Management | Beta | Planned | High | v1.2.0 | None | Core functionality |
| MobileAppCategories | App Management | Beta | Planned | Medium | v1.3.0 | None | |
| MobileAppConfigurations | App Management | Beta | Planned | Medium | v1.3.0 | MobileApps | |
| MobileAppRelationships | App Management | Beta | Planned | Medium | v1.3.0 | MobileApps | |
| MobileAppCatalogPackages | App Management | Beta | Planned | Medium | v1.3.0 | None | |
| MobileAppTroubleshootingEvents | App Management | Beta | Planned | Medium | v1.3.0 | None | |
| VppTokens | App Management | Beta | Planned | Medium | v1.3.0 | None | |
| SyncMicrosoftStoreForBusinessApps | App Management | Beta | Planned | Medium | v1.3.0 | None | |
| DeviceAppManagementTasks | App Management | Beta | Planned | Medium | v1.3.0 | None | |

### Intune App Protection Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| ManagedAppPolicies | App Protection | Beta | Planned | High | v1.2.0 | None | |
| ManagedAppRegistrations | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| ManagedAppStatuses | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| AndroidManagedAppProtections | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| IosManagedAppProtections | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsManagedAppProtections | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| DefaultManagedAppProtections | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| TargetedManagedAppConfigurations | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsInformationProtectionPolicies | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsInformationProtectionDeviceRegistrations | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsInformationProtectionWipeActions | App Protection | Beta | Planned | Medium | v1.3.0 | None | |
| MdmWindowsInformationProtectionPolicies | App Protection | Beta | Planned | Medium | v1.3.0 | None | |

### Intune eBook Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| ManagedEBooks | eBooks | Beta | Planned | Low | v1.4.0 | None | |
| ManagedEBookCategories | eBooks | Beta | Planned | Low | v1.4.0 | None | |

### Intune Code Signing Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| EnterpriseCodeSigningCertificates | Code Signing | Beta | Planned | Medium | v1.3.0 | None | |
| SymantecCodeSigningCertificate | Code Signing | Beta | Planned | Low | v1.4.0 | None | |
| WdacSupplementalPolicies | Code Signing | Beta | Planned | Medium | v1.3.0 | None | |

### Intune iOS App Provisioning Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| IosLobAppProvisioningConfigurations | iOS App Provisioning | Beta | Planned | Medium | v1.3.0 | None | |
| WindowsManagementApp | App Management | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Group Policy Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| GroupPolicyCategories | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyConfigurations | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyDefinitionFiles | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyDefinitions | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyMigrationReports | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyObjectFiles | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |
| GroupPolicyUploadedDefinitionFiles | Group Policy | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Notification and Terms Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| NotificationMessageTemplates | Notifications | Beta | Planned | Medium | v1.3.0 | None | |
| SendCustomNotificationToCompanyPortal | Notifications | Beta | Planned | Medium | v1.3.0 | None | |
| TermsAndConditions | Compliance | Beta | Planned | Medium | v1.3.0 | None | |
| IntuneBrandingProfiles | Branding | Beta | Planned | Medium | v1.3.0 | None | |

### Intune User Experience Analytics Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| UserExperienceAnalyticsOverview | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsCategories | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsBaselines | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsDevicePerformance | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsDeviceScores | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsDeviceStartupHistory | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsDeviceStartupProcesses | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsAppHealthOverview | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsAppHealthApplicationPerformance | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsBatteryHealthDevicePerformance | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsWorkFromAnywhereMetrics | Analytics | Beta | Planned | Low | v1.4.0 | None | |
| UserExperienceAnalyticsScoreHistory | Analytics | Beta | Planned | Low | v1.4.0 | None | |

### Intune Remote Management Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| RemoteAssistancePartners | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| RemoteAssistanceSettings | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| RemoteActionAudits | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| MicrosoftTunnelConfigurations | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| MicrosoftTunnelSites | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| MicrosoftTunnelHealthThresholds | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |
| VirtualEndpoint | Remote Management | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Security and Threat Defense Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| MobileThreatDefenseConnectors | Security | Beta | Planned | Medium | v1.3.0 | None | |
| AdvancedThreatProtectionOnboardingStateSummary | Security | Beta | Planned | Low | v1.4.0 | None | |
| CloudCertificationAuthority | Security | Beta | Planned | Medium | v1.3.0 | None | |
| CloudCertificationAuthorityLeafCertificate | Security | Beta | Planned | Medium | v1.3.0 | None | |
| UserPfxCertificates | Security | Beta | Planned | Medium | v1.3.0 | None | |
| EnableEndpointPrivilegeManagement | Security | Beta | Planned | Medium | v1.3.0 | None | |
| PrivilegeManagementElevations | Security | Beta | Planned | Medium | v1.3.0 | None | |
| ElevationRequests | Security | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Reports and Monitoring Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| Reports | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| AuditEvents | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| TroubleshootingEvents | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| Monitoring | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| SoftwareUpdateStatusSummary | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| IosUpdateStatuses | Reporting | Beta | Planned | Medium | v1.3.0 | None | |
| MacOSSoftwareUpdateAccountSummaries | Reporting | Beta | Planned | Medium | v1.3.0 | None | |

### Intune Integration Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| TelecomExpenseManagementPartners | Integrations | Beta | Planned | Low | v1.5.0 | None | |
| ServiceNowConnections | Integrations | Beta | Planned | Medium | v1.3.0 | None | |
| NdesConnectors | Integrations | Beta | Planned | Medium | v1.3.0 | None | |
| DomainJoinConnectors | Integrations | Beta | Planned | Medium | v1.3.0 | None | |
| EmbeddedSIMActivationCodePools | Integrations | Beta | Planned | Medium | v1.3.0 | None | |
| ZebraFotaConnector | Integrations | Beta | Planned | Low | v1.5.0 | None | |
| ZebraFotaArtifacts | Integrations | Beta | Planned | Low | v1.5.0 | None | |
| ZebraFotaDeployments | Integrations | Beta | Planned | Low | v1.5.0 | None | |
| CloudPCConnectivityIssues | Integrations | Beta | Planned | Medium | v1.4.0 | None | |

### Intune Policy Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| PolicySets | Policies | Beta | Planned | High | v1.2.0 | None | |
| Intents | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| Templates | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| TemplateSettings | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| TemplateInsights | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| ReusableSettings | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| ReusablePolicySettings | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| OperationApprovalPolicies | Policies | Beta | Planned | Medium | v1.3.0 | None | |
| OperationApprovalRequests | Policies | Beta | Planned | Medium | v1.3.0 | None | |

### Intune RBAC Resources

| Resource Name | Category | API Path | Status | Priority | Target Version | Dependencies | Notes |
|--------------|----------|----------|--------|----------|----------------|--------------|-------|
| RoleDefinitions | Intune RBAC | Beta | Completed | High | v1.0.0 | None | Core functionality |
| RoleScopeTags | Intune RBAC | Beta | Completed | High | v1.0.0 | None | Core functionality |
| RoleAssignments | Intune RBAC | Beta | Planned | High | v1.2.0 | RoleDefinitions, RoleScopeTags | Core functionality |

## Feature Request Process

If you need a resource that is not currently on our roadmap or would like to see a particular resource prioritized, please open an issue on our [GitHub repository](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues) with the "feature request" label.

## Release Schedule

We aim to release new versions on the following schedule:

- **Minor Releases (v1.x.0)**: Every 4-6 weeks
- **Patch Releases (v1.1.x)**: As needed for bug fixes
- **Major Releases (vX.0.0)**: When significant changes or breaking changes are introduced

Please note that this schedule is subject to change based on the scope of features, bug fixes, and community feedback.
