package provider

import (
	"context"
	// Graph Beta - Intune datasources
	graphBetaDeviceAndAppManagementApplicationCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/application_category"
	graphBetaDeviceAndAppManagementMobileApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app"
	graphBetaDeviceManagementAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/assignment_filter"
	graphBetaDeviceManagementDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/device_category"
	graphBetaDeviceManagementLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/linux_platform_script"
	graphBetaDeviceManagementReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/reuseable_policy_settings"
	graphBetaDeviceManagementRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/role_scope_tag"
	graphBetaDeviceManagementWindowsDriverUpdateInventory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_driver_update_inventory"
	graphBetaDeviceManagementWindowsDriverUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_driver_update_profile"
	graphBetaDeviceManagementWindowsFeatureUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_feature_update_profile"
	graphBetaDeviceManagementWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_platform_script"
	graphBetaDeviceManagementWindowsQualityUpdateExpeditePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_quality_update_expedite_policy"
	graphBetaDeviceManagementWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_quality_update_policy"
	graphBetaDeviceManagementWindowsRemediationScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_remediation_script"
	graphBetaDeviceManagementWindowsUpdateCatalogItem "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_update_catalog_item"
	graphBetaDeviceManagementWindowsUpdateRing "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_update_ring"

	// Graph Beta - M365 Admin datasources
	graphBetaM365AdminBrowserSite "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/m365_admin/graph_beta/browser_site"
	graphBetaM365AdminBrowserSiteList "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/m365_admin/graph_beta/browser_site_list"

	// Graph v1.0 - Intune datasources
	graphDeviceAndAppManagementCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_v1.0/cloud_pc_device_image"

	// Graph v1.0 - Directory Management datasources
	graphDirectoryManagementSubscribedSkus "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/directory_management/graph/subscribed_skus"

	// Utilities
	utilityMacOSPKGAppMetadata "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/utilities/macos_pkg_app_metadata"

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
		// Graph Beta - Device and app management
		graphBetaDeviceAndAppManagementApplicationCategory.NewApplicationCategoryDataSource,
		graphBetaDeviceAndAppManagementMobileApp.NewMobileAppDataSource,
		// Graph Beta - Device management
		graphBetaDeviceManagementAssignmentFilter.NewAssignmentFilterDataSource,
		graphBetaDeviceManagementDeviceCategory.NewDeviceCategoryDataSource,
		graphBetaDeviceManagementLinuxPlatformScript.NewLinuxPlatformScriptDataSource,
		graphBetaDeviceManagementReuseablePolicySettings.NewReuseablePolicySettingsDataSource,
		graphBetaDeviceManagementRoleScopeTag.NewRoleScopeTagDataSource,
		graphBetaDeviceManagementWindowsDriverUpdateProfile.NewWindowsDriverUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsDriverUpdateInventory.NewWindowsDriverUpdateInventoryDataSource,
		graphBetaDeviceManagementWindowsFeatureUpdateProfile.NewWindowsFeatureUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsQualityUpdatePolicy.NewWindowsQualityUpdateProfileDataSource,
		graphBetaDeviceManagementWindowsPlatformScript.NewWindowsPlatformScriptDataSource,
		graphBetaDeviceManagementWindowsRemediationScript.NewWindowsRemediationScriptDataSource,
		graphBetaDeviceManagementWindowsUpdateCatalogItem.NewWindowsUpdateCatalogItemDataSource,
		graphBetaDeviceManagementWindowsQualityUpdateExpeditePolicy.NewWindowsQualityUpdateExpeditePolicyDataSource,
		graphBetaDeviceManagementWindowsUpdateRing.NewWindowsUpdateRingDataSource,
		// Graph Beta - M365 Admin datasources
		graphBetaM365AdminBrowserSite.NewBrowserSiteDataSource,
		graphBetaM365AdminBrowserSiteList.NewBrowserSiteListDataSource,
		// Graph v1.0 - Intune Device and App Management datasources
		graphDeviceAndAppManagementCloudPcDeviceImage.NewCloudPcDeviceImageDataSource,
		// Graph v1.0 - Directory Management datasources
		graphDirectoryManagementSubscribedSkus.NewSubscribedSkusDataSource,

		// Utilities
		utilityMacOSPKGAppMetadata.NewMacOSPKGAppMetadataDataSource,

		// Add microsoft 365 provider datasources here
	}
}
