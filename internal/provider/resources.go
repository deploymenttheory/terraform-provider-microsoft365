package provider

import (
	"context"
	// Graph Beta - Device & app management
	graphBetaDeviceAndAppManagementApplicationCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/application_category"
	graphBetaDeviceAndAppManagementMacOSPKGApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/macos_pkg_app"
	graphBetaDeviceAndAppManagementAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/mobile_app_assignment"
	graphBetaDeviceAndAppManagementWinGetApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/winget_app"

	// Graph Beta - Device management
	graphBetaDeviceManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/assignment_filter"
	graphBetaDeviceManagementDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/device_category"
	graphBetaDeviceManagementDeviceEnrollmentConfiguration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/device_enrollment_configuration"
	graphBetaDeviceManagementEndpointPrivilegeManagement "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/endpoint_privilege_management"
	graphBetaDeviceManagementLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/linux_platform_script"
	graphBetaDeviceManagementMacOSPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/macos_platform_script"
	graphBetaDeviceManagementManagedDeviceCleanupRule "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/managed_device_cleanup_rule"

	graphBetaDeviceManagementReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/reuseable_policy_settings"
	graphBetaDeviceManagementRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/role_definition"
	graphBetaDeviceManagementRoleDefinitionAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/role_definition_assignment"
	graphBetaDeviceManagementRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/role_scope_tag"
	graphBetaDeviceManagementSettingsCatalog "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/settings_catalog"
	graphBetaDeviceManagementConfigurationPolicyAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/settings_catalog_assignment"
	graphBetaDeviceManagementSettingsCatalogJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/settings_catalog_json"
	graphBetaDeviceManagementSettingsCatalogTemplateJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/settings_catalog_template_json"
	graphBetaDeviceManagementWindowsDriverUpdateInventory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_driver_update_inventory"
	graphBetaDeviceManagementWindowsDriverUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_driver_update_profile"
	graphBetaDeviceManagementWindowsFeatureUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_feature_update_profile"
	graphBetaDeviceManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_platform_script"
	graphBetaDeviceManagementWindowsQualityExpeditePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_quality_update_expedite_policy"
	graphBetaDeviceManagementWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_quality_update_policy"
	graphBetaDeviceManagementWindowsRemediationScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_remediation_script"
	graphBetaDeviceManagementWindowsRemediationScriptAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_management/graph_beta/windows_remediation_script_assignment"

	// Graph Beta - Identity and Access resources

	"github.com/hashicorp/terraform-plugin-framework/resource"

	// Graph Beta - M365 Admin Centre
	graphBetaM365AdminBrowserSite "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/m365_admin/graph_beta/browser_site"
	graphBetaM365AdminBrowserSiteList "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/m365_admin/graph_beta/browser_site_list"
	graphDeviceM365AdminM365AppsInstallationOptions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/m365_admin/graph_beta/m365_apps_installation_options"

	// TODO current broken due to how the sdk builds time fields
	//graphBetaDeviceAndAppManagementWindowsUpdateRing "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/windows_update_ring"

	// Graph v1.0 - Intune resources
	graphDeviceAndAppManagementCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_device_image"
	graphDeviceAndAppManagementCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_provisioning_policy"
	graphDeviceAndAppManagementCloudPcUserSetting "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_user_setting"
)

// Resources returns a slice of functions that each return a resource.Resource.
// This function is a method of the M365Provider type and takes a context.Context as an argument.
// The returned slice is intended to hold the Microsoft 365 provider resources.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout.
//
// Returns:
//
//	[]func() resource.Resource: A slice of functions, each returning a resource.Resource.
//
// Resources returns a slice of functions that each return a resource.Resource.
func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Graph Beta - Device & app management
		graphBetaDeviceAndAppManagementApplicationCategory.NewApplicationCategoryResource,
		graphBetaDeviceAndAppManagementAppAssignment.NewMobileAppAssignmentResource,
		graphBetaDeviceAndAppManagementMacOSPKGApp.NewMacOSPKGAppResource,
		graphBetaDeviceAndAppManagementWinGetApp.NewWinGetAppResource,
		// Graph Beta - Device management
		graphBetaDeviceManagementAssignmentFilter.NewAssignmentFilterResource,
		graphBetaDeviceManagementDeviceCategory.NewDeviceCategoryResource,
		graphBetaDeviceManagementDeviceEnrollmentConfiguration.NewDeviceEnrollmentConfigurationResource,
		graphBetaDeviceManagementEndpointPrivilegeManagement.NewEndpointPrivilegeManagementResource,
		graphBetaDeviceManagementLinuxPlatformScript.NewLinuxPlatformScriptResource,
		graphBetaDeviceManagementMacOSPlatformScript.NewMacOSPlatformScriptResource,
		graphBetaDeviceManagementManagedDeviceCleanupRule.NewManagedDeviceCleanupRuleResource,
		graphBetaDeviceManagementSettingsCatalog.NewSettingsCatalogJsonResource,
		graphBetaDeviceManagementSettingsCatalogJson.NewSettingsCatalogJsonResource,
		graphBetaDeviceManagementSettingsCatalogTemplateJson.NewDeviceManagementTemplateJsonResource,
		graphBetaDeviceManagementConfigurationPolicyAssignment.NewDeviceManagementConfigurationPolicyAssignmentResource,
		graphBetaDeviceManagementReuseablePolicySettings.NewReuseablePolicySettingsResource,
		graphBetaDeviceManagementRoleDefinition.NewRoleDefinitionResource,
		graphBetaDeviceManagementRoleDefinitionAssignment.NewRoleDefinitionAssignmentResource,
		graphBetaDeviceManagementRoleScopeTag.NewRoleScopeTagResource,
		graphBetaDeviceManagementWindowsDriverUpdateProfile.NewWindowsDriverUpdateProfileResource,
		graphBetaDeviceManagementWindowsDriverUpdateInventory.NewWindowsDriverUpdateInventoryResource,
		graphBetaDeviceManagementWindowsFeatureUpdateProfile.NewWindowsFeatureUpdateProfileResource,
		graphBetaDeviceManagementWindowsPlatformScript.NewWindowsPlatformScriptResource,
		graphBetaDeviceManagementWindowsRemediationScript.NewDeviceHealthScriptResource,
		graphBetaDeviceManagementWindowsRemediationScriptAssignment.NewDeviceHealthScriptAssignmentResource,
		graphBetaDeviceManagementWindowsQualityExpeditePolicy.NewWindowsQualityUpdateExpeditePolicyResource,
		graphBetaDeviceManagementWindowsQualityUpdatePolicy.NewWindowsQualityUpdatePolicyResource,
		//graphBetaDeviceAndAppManagementWindowsUpdateRing.NewWindowsUpdateRingResource,

		// Graph Beta - Identity and Access resources
		//graphBetaIdentityAndAccessConditionalAccessPolicy.NewConditionalAccessPolicyResource,

		// Graph Beta - M365 Admin Centre
		graphBetaM365AdminBrowserSite.NewBrowserSiteResource,
		graphBetaM365AdminBrowserSiteList.NewBrowserSiteListResource,

		// Graph v1.0 - Intune resources
		graphDeviceAndAppManagementCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource,
		graphDeviceAndAppManagementCloudPcUserSetting.NewCloudPcUserSettingResource,
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageResource,

		// Graph v1.0 - M365 Admin Centre
		graphDeviceM365AdminM365AppsInstallationOptions.NewM365AppsInstallationOptionsResource,
		// Add microsoft 365 provider resources here
	}
}
