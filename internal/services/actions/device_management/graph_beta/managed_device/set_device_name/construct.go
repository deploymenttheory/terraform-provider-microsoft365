package graphBetaSetDeviceNameManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceSetName) *devicemanagement.ManagedDevicesItemSetDeviceNamePostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemSetDeviceNamePostRequestBody()

	convert.FrameworkToGraphString(device.DeviceName, requestBody.SetDeviceName)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device set device name request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceSetName) *devicemanagement.ComanagedDevicesItemSetDeviceNamePostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemSetDeviceNamePostRequestBody()

	convert.FrameworkToGraphString(device.DeviceName, requestBody.SetDeviceName)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device set device name request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

