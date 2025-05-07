package provider

import (
	"context"
	// Graph Beta - Intune datasources
	graphBetaDeviceAndAppManagementMacOSPKGApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/macos_pkg_app"
	graphBetaDeviceManagementApplicationCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/application_category"
	graphBetaDeviceManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/assignment_filter"
	graphBetaDeviceManagementDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/device_category"
	graphBetaDeviceManagementLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/linux_platform_script"
	graphBetaDeviceManagementReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/reuseable_policy_settings"
	graphBetaDeviceManagementRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/role_scope_tag"
	graphBetaDeviceManagementWindowsDriverUpdateInventory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_driver_update_inventory"
	graphBetaDeviceManagementWindowsDriverUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_driver_update_profile"
	graphBetaDeviceManagementWindowsFeatureUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_feature_update_profile"
	graphBetaDeviceManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_platform_script"
	graphBetaDeviceManagementWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_quality_update_policy"
	graphBetaDeviceManagementWindowsUpdateCatalogItem "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_management/graph_beta/windows_update_catalog_item"

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
		graphBetaDeviceManagementApplicationCategory.NewApplicationCategoryDataSource,
		graphBetaDeviceManagementAssignmentFilter.NewAssignmentFilterDataSource,
		graphBetaDeviceManagementDeviceCategory.NewDeviceCategoryDataSource,
		graphBetaDeviceManagementLinuxPlatformScript.NewLinuxPlatformScriptDataSource,
		graphBetaDeviceAndAppManagementMacOSPKGApp.NewMacOSPKGAppDataSource,
		graphBetaDeviceManagementReuseablePolicySettings.NewReuseablePolicySettingsDataSource,
		graphBetaDeviceManagementRoleScopeTag.NewRoleScopeTagDataSource,
		graphBetaDeviceManagementWindowsDriverUpdateProfile.NewWindowsDriverUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsDriverUpdateInventory.NewWindowsDriverUpdateInventoryDataSource,
		graphBetaDeviceManagementWindowsFeatureUpdateProfile.NewWindowsFeatureUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsQualityUpdatePolicy.NewWindowsQualityUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsPlatformScript.NewWindowsPlatformScriptDataSource,
		graphBetaDeviceManagementWindowsUpdateCatalogItem.NewWindowsUpdateCatalogItemDataSource,
		// Graph v1.0 - Intune datasources
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageDataSource,

		// Add microsoft 365 provider datasources here
	}
}
