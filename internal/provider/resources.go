package provider

import (
	"context"
	// Graph Beta - Intune resources
	graphBetaDeviceAndAppManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/assignment_filter"
	graphBetaDeviceAndAppManagementBrowserSite "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/browser_site"
	graphBetaDeviceAndAppManagementBrowserSiteList "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/browser_site_list"
	graphBetaDeviceAndAppManagementEndpointPrivilegeManagement "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/endpoint_privilege_management"
	graphBetaDeviceAndAppManagementLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/linux_platform_script"
	graphBetaDeviceAndAppManagementM365AppsInstallationOptions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/m365_apps_installation_options"
	graphBetaDeviceAndAppManagementmacOSPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/macos_platform_script"
	graphBetaDeviceAndAppManagementReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/reuseable_policy_settings"
	graphBetaDeviceAndAppManagementRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/role_definition"
	graphBetaDeviceAndAppManagementRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/role_scope_tag"
	graphBetaDeviceAndAppManagementSettingsCatalog "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/settings_catalog"
	graphBetaDeviceAndAppManagementSettingsCatalogTemplate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/settings_catalog_template"
	graphBetaDeviceAndAppManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/windows_platform_script"
	graphBetaDeviceAndAppManagementWinGetApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/winget_app"

	// Graph v1.0 - Intune resources
	graphDeviceAndAppManagementCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_device_image"
	graphDeviceAndAppManagementCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_provisioning_policy"
	graphDeviceAndAppManagementCloudPcUserSetting "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_v1.0/cloud_pc_user_setting"

	// Graph Beta - Identity and Access resources
	graphBetaIdentityAndAccessConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/identity_and_access/graph_beta/conditional_access_policy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
		// Graph Beta - Intune resources
		graphBetaDeviceAndAppManagementApplications.NewApplicationsResource,
		graphBetaDeviceAndAppManagementAssignmentFilter.NewAssignmentFilterResource,
		graphBetaDeviceAndAppManagementBrowserSite.NewBrowserSiteResource,
		graphBetaDeviceAndAppManagementBrowserSiteList.NewBrowserSiteListResource,
		graphBetaDeviceAndAppManagementEndpointPrivilegeManagement.NewEndpointPrivilegeManagementResource,
		graphBetaDeviceAndAppManagementLinuxPlatformScript.NewLinuxPlatformScriptResource,
		graphBetaDeviceAndAppManagementmacOSPlatformScript.NewMacOSPlatformScriptResource,
		graphBetaDeviceAndAppManagementM365AppsInstallationOptions.NewM365AppsInstallationOptionsResource,
		graphBetaDeviceAndAppManagementSettingsCatalog.NewSettingsCatalogResource,
		graphBetaDeviceAndAppManagementSettingsCatalogTemplate.NewDeviceManagementTemplateResource,
		graphBetaDeviceAndAppManagementReuseablePolicySettings.NewReuseablePolicySettingsResource,
		graphBetaDeviceAndAppManagementRoleDefinition.NewRoleDefinitionResource,
		graphBetaDeviceAndAppManagementRoleScopeTag.NewRoleScopeTagResource,
		graphBetaDeviceAndAppManagementWindowsPlatformScript.NewWindowsPlatformScriptResource,
		graphBetaDeviceAndAppManagementWinGetApp.NewWinGetAppResource,
		// Graph Beta - Identity and Access resources
		graphBetaIdentityAndAccessConditionalAccessPolicy.NewConditionalAccessPolicyResource,
		// Graph v1.0 - Intune resources
		graphDeviceAndAppManagementCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource,
		graphDeviceAndAppManagementCloudPcUserSetting.NewCloudPcUserSettingResource,
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageResource,
		// Add microsoft 365 provider resources here
	}
}
