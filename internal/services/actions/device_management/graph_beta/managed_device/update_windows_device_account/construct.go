package graphBetaUpdateWindowsDeviceAccount

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructManagedDeviceRequest builds the request body for updating managed device account
func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceAccount) *devicemanagement.ManagedDevicesItemUpdateWindowsDeviceAccountPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemUpdateWindowsDeviceAccountPostRequestBody()

	actionParameter := models.NewUpdateWindowsDeviceAccountActionParameter()

	deviceAccount := models.NewWindowsDeviceAccount()
	convert.FrameworkToGraphString(device.Password, deviceAccount.SetPassword)
	actionParameter.SetDeviceAccount(deviceAccount)

	convert.FrameworkToGraphBool(device.PasswordRotationEnabled, actionParameter.SetPasswordRotationEnabled)
	convert.FrameworkToGraphBool(device.CalendarSyncEnabled, actionParameter.SetCalendarSyncEnabled)
	convert.FrameworkToGraphString(device.DeviceAccountEmail, actionParameter.SetDeviceAccountEmail)

	convert.FrameworkToGraphString(device.ExchangeServer, actionParameter.SetExchangeServer)
	convert.FrameworkToGraphString(device.SessionInitiationProtocolAddress, actionParameter.SetSessionInitiationProtocalAddress)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device account update request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	requestBody.SetUpdateWindowsDeviceAccountActionParameter(actionParameter)
	return requestBody
}

// constructComanagedDeviceRequest builds the request body for updating co-managed device account
func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceAccount) *devicemanagement.ComanagedDevicesItemUpdateWindowsDeviceAccountPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemUpdateWindowsDeviceAccountPostRequestBody()

	actionParameter := models.NewUpdateWindowsDeviceAccountActionParameter()

	deviceAccount := models.NewWindowsDeviceAccount()
	convert.FrameworkToGraphString(device.Password, deviceAccount.SetPassword)
	actionParameter.SetDeviceAccount(deviceAccount)

	convert.FrameworkToGraphBool(device.PasswordRotationEnabled, actionParameter.SetPasswordRotationEnabled)
	convert.FrameworkToGraphBool(device.CalendarSyncEnabled, actionParameter.SetCalendarSyncEnabled)
	convert.FrameworkToGraphString(device.DeviceAccountEmail, actionParameter.SetDeviceAccountEmail)

	convert.FrameworkToGraphString(device.ExchangeServer, actionParameter.SetExchangeServer)
	convert.FrameworkToGraphString(device.SessionInitiationProtocolAddress, actionParameter.SetSessionInitiationProtocalAddress)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device account update request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	requestBody.SetUpdateWindowsDeviceAccountActionParameter(actionParameter)
	return requestBody
}
