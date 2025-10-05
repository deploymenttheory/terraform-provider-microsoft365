package graphBetaWindowsBackupAndRestore

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GraphServiceClient object to a Terraform state.
// Returns true if mapping was successful, false if the resource is not a Windows Restore configuration.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsBackupAndRestoreResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) bool {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return false
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	windowsRestoreConfig, ok := remoteResource.(*graphmodels.WindowsRestoreDeviceEnrollmentConfiguration)
	if !ok {
		tflog.Warn(ctx, "Remote resource is not a Windows Restore device enrollment configuration")
		return false
	}
	data.State = convert.GraphToFrameworkEnum(windowsRestoreConfig.GetState())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return true
}
