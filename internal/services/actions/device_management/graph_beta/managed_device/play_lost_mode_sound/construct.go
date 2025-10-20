package graphBetaPlayLostModeSoundManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDevicePlaySound) *devicemanagement.ManagedDevicesItemPlayLostModeSoundPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemPlayLostModeSoundPostRequestBody()

	convert.FrameworkToGraphString(device.DurationInMinutes, requestBody.SetDurationInMinutes)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device play lost mode sound request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDevicePlaySound) *devicemanagement.ComanagedDevicesItemPlayLostModeSoundPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemPlayLostModeSoundPostRequestBody()

	convert.FrameworkToGraphString(device.DurationInMinutes, requestBody.SetDurationInMinutes)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device play lost mode sound request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
