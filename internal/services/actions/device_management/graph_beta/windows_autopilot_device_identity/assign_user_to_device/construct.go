package graphBetaAssignUserToDevice

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructRequest constructs the request body for the assignUserToDevice action
func constructRequest(ctx context.Context, data *AssignUserToDeviceActionModel) (*devicemanagement.WindowsAutopilotDeviceIdentitiesItemAssignUserToDevicePostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing assignUserToDevice request for device ID: %s", data.WindowsAutopilotDeviceIdentityID.ValueString()))

	requestBody := devicemanagement.NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDevicePostRequestBody()

	userPrincipalName := data.UserPrincipalName.ValueString()
	addressableUserName := data.AddressableUserName.ValueString()

	requestBody.SetUserPrincipalName(&userPrincipalName)
	requestBody.SetAddressableUserName(&addressableUserName)

	tflog.Debug(ctx, fmt.Sprintf("Request body created with UserPrincipalName: %s, AddressableUserName: %s",
		userPrincipalName, addressableUserName))

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignUserToDevice request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignUserToDevice request", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing assignUserToDevice request")
	return requestBody, nil
}
