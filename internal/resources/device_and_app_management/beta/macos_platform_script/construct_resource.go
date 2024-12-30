// Main entry point to construct the intune windows device management script resource for the Terraform provider.
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows device management script resource for the Terraform provider.
func constructResource(ctx context.Context, data *MacOSPlatformScriptResourceModel) (graphmodels.DeviceShellScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceShellScript()

	construct.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	construct.SetStringProperty(data.Description, requestBody.SetDescription)
	construct.SetBytesProperty(data.ScriptContent, requestBody.SetScriptContent)

	if err := construct.SetEnumProperty(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	construct.SetStringProperty(data.FileName, requestBody.SetFileName)

	if err := construct.SetStringList(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	construct.SetBoolProperty(data.BlockExecutionNotifications, requestBody.SetBlockExecutionNotifications)

	if !data.ExecutionFrequency.IsNull() {
		isoDuration, err := serialization.ParseISODuration(data.ExecutionFrequency.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing execution frequency: %v", err)
		}
		requestBody.SetExecutionFrequency(isoDuration)
	}

	construct.SetInt32Property(data.RetryCount, requestBody.SetRetryCount)

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
