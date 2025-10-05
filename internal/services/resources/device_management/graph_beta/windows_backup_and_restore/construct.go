package graphBetaWindowsBackupAndRestore

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the Graph API request body for Windows Backup and Restore configuration
func constructResource(ctx context.Context, data *WindowsBackupAndRestoreResourceModel) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing Windows Backup and Restore configuration request body")

	config := graphmodels.NewWindowsRestoreDeviceEnrollmentConfiguration()

	err := convert.FrameworkToGraphEnum(data.State, graphmodels.ParseEnablement, config.SetState)
	if err != nil {
		return nil, fmt.Errorf("error setting State: %v", err)
	}

	tflog.Debug(ctx, "Successfully constructed Windows Backup and Restore configuration request body", map[string]any{
		"state": data.State.ValueString(),
	})

	return config, nil
}
