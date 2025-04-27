package provider

import (
	"context"
	// Graph Beta - Intune datasources
	graphBetaDeviceAndAppManagementApplicationCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/application_category"
	graphBetaDeviceAndAppManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/assignment_filter"
	graphBetaDeviceAndAppManagementDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/device_category"
	graphBetaDeviceAndAppManagementLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/linux_platform_script"
	graphBetaDeviceAndAppManagementMacOSPKGApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/macos_pkg_app"
	graphBetaDeviceAndAppManagementReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/reuseable_policy_settings"
	graphBetaDeviceAndAppManagementRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/role_scope_tag"
	graphBetaDeviceAndAppManagementWindowsDriverUpdateInventory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_driver_update_inventory"
	graphBetaDeviceAndAppManagementWindowsDriverUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_driver_update_profile"
	graphBetaDeviceAndAppManagementWindowsFeatureUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_feature_update_profile"
	graphBetaDeviceAndAppManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_platform_script"
	graphBetaDeviceAndAppManagementWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_quality_update_policy"
	graphBetaDeviceAndAppManagementWindowsUpdateCatalogItem "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_update_catalog_item"

	// Graph v1.0 - Intune datasources
	graphDeviceAndAppManagementCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_v1.0/cloud_pc_device_image"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns a slice of functions that each return a datasource.DataSource.
// This function is a method of the M365Provider type and takes a context.Context as an argument.
// The returned slice is intended to hold the Microsoft 365 provider datasources.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout.
//
// Returns:
//
//	[]func() datasource.DataSource: A slice of functions, each returning a datasource.DataSource.
func (p *M365Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Graph Beta - Intune datasources
		graphBetaDeviceAndAppManagementApplicationCategory.NewApplicationCategoryDataSource,
		graphBetaDeviceAndAppManagementAssignmentFilter.NewAssignmentFilterDataSource,
		graphBetaDeviceAndAppManagementDeviceCategory.NewDeviceCategoryDataSource,
		graphBetaDeviceAndAppManagementLinuxPlatformScript.NewLinuxPlatformScriptDataSource,
		graphBetaDeviceAndAppManagementMacOSPKGApp.NewMacOSPKGAppDataSource,
		graphBetaDeviceAndAppManagementReuseablePolicySettings.NewReuseablePolicySettingsDataSource,
		graphBetaDeviceAndAppManagementRoleScopeTag.NewRoleScopeTagDataSource,
		graphBetaDeviceAndAppManagementWindowsDriverUpdateProfile.NewWindowsDriverUpdateProfileDataSource,
		graphBetaDeviceAndAppManagementWindowsDriverUpdateInventory.NewWindowsDriverUpdateInventoryDataSource,
		graphBetaDeviceAndAppManagementWindowsFeatureUpdateProfile.NewWindowsFeatureUpdateProfileDataSource,
		graphBetaDeviceAndAppManagementWindowsQualityUpdatePolicy.NewWindowsQualityUpdateProfileDataSource,
		graphBetaDeviceAndAppManagementWindowsPlatformScript.NewWindowsPlatformScriptDataSource,
		graphBetaDeviceAndAppManagementWindowsUpdateCatalogItem.NewWindowsUpdateCatalogItemDataSource,
		// Graph v1.0 - Intune datasources
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageDataSource,

		// Add microsoft 365 provider datasources here
	}
}
