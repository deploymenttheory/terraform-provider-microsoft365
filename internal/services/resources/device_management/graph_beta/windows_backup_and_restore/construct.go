package graphBetaWindowsBackupAndRestore

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the Graph API request body for Windows Backup and Restore configuration
func constructResource(ctx context.Context, data *WindowsBackupAndRestoreResourceModel, isCreate bool) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing Windows Backup and Restore configuration request body")

	// Create the Windows Restore device enrollment configuration
	config := graphmodels.NewWindowsRestoreDeviceEnrollmentConfiguration()

	// Set basic properties using conversion helpers
	convert.FrameworkToGraphString(data.DisplayName, config.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, config.SetDescription)
	convert.FrameworkToGraphInt32(data.Priority, config.SetPriority)

	// Set role scope tag IDs using conversion helper
	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, config.SetRoleScopeTagIds)
		if err != nil {
			return nil, err
		}
	}

	// Set the device enrollment configuration type
	deviceEnrollmentConfigurationType := graphmodels.WINDOWSRESTORE_DEVICEENROLLMENTCONFIGURATIONTYPE
	config.SetDeviceEnrollmentConfigurationType(&deviceEnrollmentConfigurationType)

	// Set the state using conversion helper
	if !data.State.IsNull() && !data.State.IsUnknown() {
		err := convert.FrameworkToGraphEnum(data.State, graphmodels.ParseEnablement, config.SetState)
		if err != nil {
			return nil, err
		}
	}

	// Set OData type
	odataType := "#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration"
	config.SetOdataType(&odataType)

	tflog.Debug(ctx, "Successfully constructed Windows Backup and Restore configuration request body", map[string]any{
		"displayName": data.DisplayName.ValueString(),
		"state":       data.State.ValueString(),
	})

	return config, nil
}
