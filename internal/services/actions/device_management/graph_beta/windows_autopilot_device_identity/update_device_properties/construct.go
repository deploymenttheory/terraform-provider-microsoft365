package graphBetaUpdateDeviceProperties

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructRequest constructs the request body for the updateDeviceProperties action
func constructRequest(ctx context.Context, data *UpdateDevicePropertiesActionModel) (*devicemanagement.WindowsAutopilotDeviceIdentitiesItemUpdateDevicePropertiesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s request for device ID: %s", ActionName, data.WindowsAutopilotDeviceIdentityID.ValueString()))

	requestBody := devicemanagement.NewWindowsAutopilotDeviceIdentitiesItemUpdateDevicePropertiesPostRequestBody()

	convert.FrameworkToGraphString(data.UserPrincipalName, requestBody.SetUserPrincipalName)
	convert.FrameworkToGraphString(data.AddressableUserName, requestBody.SetAddressableUserName)
	convert.FrameworkToGraphString(data.GroupTag, requestBody.SetGroupTag)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for action %s", ActionName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s request for %s", ActionName, data.WindowsAutopilotDeviceIdentityID.ValueString()))
	return requestBody, nil
}
