// Main entry point to construct the intune windows device management script resource for the Terraform provider.
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows device management script resource for the Terraform provider.
func constructResource(ctx context.Context, data *MacOSPlatformScriptResourceModel) (graphmodels.DeviceShellScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceShellScript()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetBytesProperty(data.ScriptContent, requestBody.SetScriptContent)

	if err := constructors.SetEnumProperty(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	constructors.SetStringProperty(data.FileName, requestBody.SetFileName)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	constructors.SetBoolProperty(data.BlockExecutionNotifications, requestBody.SetBlockExecutionNotifications)

	if err := constructors.SetISODurationProperty(data.ExecutionFrequency, requestBody.SetExecutionFrequency); err != nil {
		return nil, fmt.Errorf("error setting execution frequency: %v", err)
	}

	constructors.SetInt32Property(data.RetryCount, requestBody.SetRetryCount)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
