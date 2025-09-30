package graphBetaUpdateDeviceProperties

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructRequest constructs the request body for the updateDeviceProperties action
func constructRequest(ctx context.Context, data *UpdateDevicePropertiesActionModel) (*devicemanagement.WindowsAutopilotDeviceIdentitiesItemUpdateDevicePropertiesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing updateDeviceProperties request for device ID: %s", data.WindowsAutopilotDeviceIdentityID.ValueString()))

	// Create the request body for the updateDeviceProperties action
	requestBody := devicemanagement.NewWindowsAutopilotDeviceIdentitiesItemUpdateDevicePropertiesPostRequestBody()

	// Set user principal name if provided
	if !data.UserPrincipalName.IsNull() && !data.UserPrincipalName.IsUnknown() {
		userPrincipalName := data.UserPrincipalName.ValueString()
		requestBody.SetUserPrincipalName(&userPrincipalName)
		tflog.Debug(ctx, fmt.Sprintf("Setting UserPrincipalName: %s", userPrincipalName))
	}

	// Set addressable user name if provided
	if !data.AddressableUserName.IsNull() && !data.AddressableUserName.IsUnknown() {
		addressableUserName := data.AddressableUserName.ValueString()
		requestBody.SetAddressableUserName(&addressableUserName)
		tflog.Debug(ctx, fmt.Sprintf("Setting AddressableUserName: %s", addressableUserName))
	}

	// Set group tag if provided
	if !data.GroupTag.IsNull() && !data.GroupTag.IsUnknown() {
		groupTag := data.GroupTag.ValueString()
		requestBody.SetGroupTag(&groupTag)
		tflog.Debug(ctx, fmt.Sprintf("Setting GroupTag: %s", groupTag))
	}

	// Set display name if provided
	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
		tflog.Debug(ctx, fmt.Sprintf("Setting DisplayName: %s", displayName))
	}

	// Debug log the constructed object
	if err := constructors.DebugLogGraphObject(ctx, "Constructed updateDeviceProperties request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log updateDeviceProperties request", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing updateDeviceProperties request")
	return requestBody, nil
}