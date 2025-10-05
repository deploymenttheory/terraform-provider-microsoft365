package graphBetaWindowsBackupAndRestore

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// resolveWindowsRestoreDeviceEnrollmentConfigurationID queries the Microsoft Graph API
// to find the Windows Restore Device Enrollment Configuration and returns its ID.
// This function looks for configurations with @odata.type "#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration"
// and returns the ID of the first matching configuration.
//
// Parameters:
//   - ctx: The context for the request
//   - client: The Microsoft Graph Beta SDK client
//   - resp: The Terraform response for error handling
//   - operation: The operation name for error context
//   - writePermissions: Required permissions for error handling
//
// Returns:
//   - string: The ID of the Windows Restore Device Enrollment Configuration (empty if error)
//
// Example returned ID: "54fac284-7866-43e5-860a-9c8e10fa3d7d_WindowsRestore"
func resolveWindowsRestoreDeviceEnrollmentConfigurationID(ctx context.Context, client *devicemanagement.DeviceManagementRequestBuilder, resp any, operation string, writePermissions []string) string {
	tflog.Debug(ctx, "Starting resolution of Windows Restore Device Enrollment Configuration ID")

	configurations, err := client.
		DeviceEnrollmentConfigurations().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, writePermissions)
		return ""
	}

	if configurations == nil || configurations.GetValue() == nil {
		tflog.Warn(ctx, "No device enrollment configurations found")
		return ""
	}

	for _, config := range configurations.GetValue() {
		if config == nil {
			continue
		}

		if windowsRestoreConfig, ok := config.(*graphmodels.WindowsRestoreDeviceEnrollmentConfiguration); ok {
			configID := windowsRestoreConfig.GetId()
			if configID != nil && *configID != "" {
				tflog.Debug(ctx, "Found Windows Restore Device Enrollment Configuration", map[string]any{
					"id":          *configID,
					"displayName": helpers.GetStringValue(windowsRestoreConfig.GetDisplayName()),
					"description": helpers.GetStringValue(windowsRestoreConfig.GetDescription()),
				})
				return *configID
			}
		}
	}

	tflog.Warn(ctx, "Windows Restore Device Enrollment Configuration not found")
	return ""
}
