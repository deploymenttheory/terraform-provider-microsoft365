package graphBetaEnableLostModeManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceLostMode) *devicemanagement.ManagedDevicesItemEnableLostModePostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemEnableLostModePostRequestBody()

	convert.FrameworkToGraphString(device.Message, requestBody.SetMessage)
	convert.FrameworkToGraphString(device.PhoneNumber, requestBody.SetPhoneNumber)
	convert.FrameworkToGraphString(device.Footer, requestBody.SetFooter)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device enable lost mode request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceLostMode) *devicemanagement.ComanagedDevicesItemEnableLostModePostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemEnableLostModePostRequestBody()

	convert.FrameworkToGraphString(device.Message, requestBody.SetMessage)
	convert.FrameworkToGraphString(device.PhoneNumber, requestBody.SetPhoneNumber)
	convert.FrameworkToGraphString(device.Footer, requestBody.SetFooter)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device enable lost mode request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
