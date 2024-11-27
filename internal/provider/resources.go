package provider

import (
	"context"

	graphBetaDeviceAndAppManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/assignment_filter"
	graphBetaDeviceAndAppManagementBrowserSite "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/browser_site"
	graphBetaDeviceAndAppManagementBrowserSiteList "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/browser_site_list"
	graphBetaDeviceAndAppManagementM365AppsInstallationOptions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/m365_apps_installation_options"
	graphBetaDeviceAndAppManagementmacOSPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/macos_platform_script"
	graphBetaDeviceAndAppManagementMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/mobile_app_assignment"
	graphBetaDeviceAndAppManagementRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/role_definition"
	graphBetaDeviceAndAppManagementSettingsCatalog "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/settings_catalog"
	graphBetaDeviceAndAppManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/windows_platform_script"
	graphBetaDeviceAndAppManagementWinGetApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/winget_app"
	graphDeviceAndAppManagementCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_device_image"
	graphDeviceAndAppManagementCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_provisioning_policy"
	graphDeviceAndAppManagementCloudPcUserSetting "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_user_setting"
	graphDeviceAndAppManagementRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/role_definition"
	graphBetaIdentityAndAccessConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/identity_and_access/beta/conditional_access_policy"
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
		graphBetaDeviceAndAppManagementAssignmentFilter.NewAssignmentFilterResource,
		graphBetaDeviceAndAppManagementBrowserSite.NewBrowserSiteResource,
		graphBetaDeviceAndAppManagementBrowserSiteList.NewBrowserSiteListResource,
		graphBetaDeviceAndAppManagementmacOSPlatformScript.NewDeviceShellScriptResource,
		graphBetaDeviceAndAppManagementM365AppsInstallationOptions.NewM365AppsInstallationOptionsResource,
		graphBetaDeviceAndAppManagementMobileAppAssignment.NewMobileAppAssignmentResource,
		graphBetaDeviceAndAppManagementSettingsCatalog.NewSettingsCatalogResource,
		graphBetaDeviceAndAppManagementRoleDefinition.NewRoleDefinitionResource,
		graphBetaDeviceAndAppManagementWindowsPlatformScript.NewWindowsPlatformScriptResource,
		graphBetaDeviceAndAppManagementWinGetApp.NewWinGetAppResource,
		graphBetaIdentityAndAccessConditionalAccessPolicy.NewConditionalAccessPolicyResource,
		graphDeviceAndAppManagementCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource,
		graphDeviceAndAppManagementCloudPcUserSetting.NewCloudPcUserSettingResource,
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageResource,
		graphDeviceAndAppManagementRoleDefinition.NewRoleDefinitionResource,
		// Add microsoft 365 provider resources here
	}
}
