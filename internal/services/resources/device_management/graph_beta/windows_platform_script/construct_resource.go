package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows device management script resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsPlatformScriptResourceModel) (graphmodels.DeviceManagementScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementScript()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBytes(data.ScriptContent, requestBody.SetScriptContent)

	if err := convert.FrameworkToGraphEnum(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("failed to set runAsAccount: %v", err)
	}

	convert.FrameworkToGraphBool(data.EnforceSignatureCheck, requestBody.SetEnforceSignatureCheck)
	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	convert.FrameworkToGraphBool(data.RunAs32Bit, requestBody.SetRunAs32Bit)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
