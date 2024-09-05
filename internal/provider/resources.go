package provider

import (
	"context"

	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/assignment_filter"
	graphBetaDeviceManagementScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/device_management_script"
	graphBetaM365AppsInstallationOptions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/m365_apps_installation_options"
	graphCloudPcDeviceImage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_device_image"
	graphCloudPcProvisioningPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_provisioning_policy"
	graphCloudPcUserSetting "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/v1.0/cloud_pc_user_setting"
	graphBetaConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/identity_and_access/beta/conditional_access_policy"

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
		graphBetaAssignmentFilter.NewAssignmentFilterResource,
		graphBetaDeviceManagementScript.NewDeviceManagementScriptResource,
		graphBetaConditionalAccessPolicy.NewConditionalAccessPolicyResource,
		graphBetaM365AppsInstallationOptions.NewM365AppsInstallationOptionsResource,
		graphCloudPcProvisioningPolicy.NewCloudPcProvisioningPolicyResource,
		graphCloudPcUserSetting.NewCloudPcUserSettingResource,
		graphCloudPcDeviceImage.NewCloudPcDeviceImageResource,
		// Add microsoft 365 provider resources here
	}
}
