// Main entry point to construct the intune windows device management script resource for the Terraform provider.
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows device management script resource for the Terraform provider.
func constructResource(ctx context.Context, data *MacOSPlatformScriptResourceModel) (graphmodels.DeviceShellScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceShellScript()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBytes(data.ScriptContent, requestBody.SetScriptContent)

	if err := convert.FrameworkToGraphEnum(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	convert.FrameworkToGraphBool(data.BlockExecutionNotifications, requestBody.SetBlockExecutionNotifications)

	if err := convert.FrameworkToGraphISODuration(data.ExecutionFrequency, requestBody.SetExecutionFrequency); err != nil {
		return nil, fmt.Errorf("error setting execution frequency: %v", err)
	}

	convert.FrameworkToGraphInt32(data.RetryCount, requestBody.SetRetryCount)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
