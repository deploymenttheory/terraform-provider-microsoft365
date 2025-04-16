package provider

import (
	"context"
	// Graph Beta - Intune resources
	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/assignment_filter"
	graphBetaDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/device_category"
	graphBetaLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/linux_platform_script"
	graphBetaMacOSPKGApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/macos_pkg_app"
	graphBetaReuseablePolicySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/reuseable_policy_settings"
	graphBetaRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/role_scope_tag"
	graphBetaWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_beta/windows_platform_script"

	// Graph v1.0 - Intune resources
	graphCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/device_and_app_management/graph_v1.0/cloud_pc_device_image"

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
		graphBetaAssignmentFilter.NewAssignmentFilterDataSource,
		graphBetaDeviceCategory.NewDeviceCategoryDataSource,
		graphBetaLinuxPlatformScript.NewLinuxPlatformScriptDataSource,
		graphBetaMacOSPKGApp.NewMacOSPKGAppDataSource,
		graphBetaReuseablePolicySettings.NewReuseablePolicySettingsDataSource,
		graphBetaRoleScopeTag.NewRoleScopeTagDataSource,
		graphBetaWindowsPlatformScript.NewWindowsPlatformScriptDataSource,
		graphCloudPcDeviceImage.NewCloudPcDeviceImageDataSource,

		// Add microsoft 365 provider datasources here
	}
}
