package graphBetaMacOSCustomAttributeScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *DeviceCustomAttributeShellScriptResourceModel) (msgraphbetamodels.DeviceCustomAttributeShellScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))
	requestBody := msgraphbetamodels.NewDeviceCustomAttributeShellScript()

	convert.FrameworkToGraphString(data.CustomAttributeName, requestBody.SetCustomAttributeName)
	if err := convert.FrameworkToGraphEnum(data.CustomAttributeType, msgraphbetamodels.ParseDeviceCustomAttributeValueType, requestBody.SetCustomAttributeType); err != nil {
		return nil, err
	}
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBytes(data.ScriptContent, requestBody.SetScriptContent)

	if err := convert.FrameworkToGraphEnum(data.RunAsAccount, msgraphbetamodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, err
	}

	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)
	convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
