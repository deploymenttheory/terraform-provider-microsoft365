package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new DeviceHealthScript based on the resource model
func constructResource(ctx context.Context, data *DeviceHealthScriptResourceModel) (graphmodels.DeviceHealthScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceHealthScript()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphBool(data.RunAs32Bit, requestBody.SetRunAs32Bit)
	convert.FrameworkToGraphBool(data.EnforceSignatureCheck, requestBody.SetEnforceSignatureCheck)

	if err := convert.FrameworkToGraphEnum(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	if data.DeviceHealthScriptType.ValueString() != "" {
		if err := convert.FrameworkToGraphEnum(data.DeviceHealthScriptType, graphmodels.ParseDeviceHealthScriptType, requestBody.SetDeviceHealthScriptType); err != nil {
			return nil, fmt.Errorf("invalid device health script type: %s", err)
		}
	}

	convert.FrameworkToGraphBytes(data.DetectionScriptContent, requestBody.SetDetectionScriptContent)
	convert.FrameworkToGraphBytes(data.RemediationScriptContent, requestBody.SetRemediationScriptContent)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
