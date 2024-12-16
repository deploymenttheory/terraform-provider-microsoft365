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

	if !data.DisplayName.IsNull() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.ScriptContent.IsNull() {
		scriptContent := []byte(data.ScriptContent.ValueString())
		requestBody.SetScriptContent(scriptContent)
	}

	if !data.RunAsAccount.IsNull() {
		runAsAccountStr := data.RunAsAccount.ValueString()
		var runAsAccount graphmodels.RunAsAccountType
		switch runAsAccountStr {
		case "system":
			runAsAccount = graphmodels.SYSTEM_RUNASACCOUNTTYPE
		case "user":
			runAsAccount = graphmodels.USER_RUNASACCOUNTTYPE
		}
		requestBody.SetRunAsAccount(&runAsAccount)
	}

	if !data.FileName.IsNull() {
		fileName := data.FileName.ValueString()
		requestBody.SetFileName(&fileName)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			requestBody.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if !data.BlockExecutionNotifications.IsNull() {
		blockNotifications := data.BlockExecutionNotifications.ValueBool()
		requestBody.SetBlockExecutionNotifications(&blockNotifications)
	}

	if !data.ExecutionFrequency.IsNull() {
		isoDuration, err := serialization.ParseISODuration(data.ExecutionFrequency.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing execution frequency: %v", err)
		}
		requestBody.SetExecutionFrequency(isoDuration)
	}

	if !data.RetryCount.IsNull() {
		retryCount := data.RetryCount.ValueInt32()
		requestBody.SetRetryCount(&retryCount)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
