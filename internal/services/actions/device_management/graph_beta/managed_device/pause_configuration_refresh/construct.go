package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDevicePauseConfig) *devicemanagement.ManagedDevicesItemPauseConfigurationRefreshPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemPauseConfigurationRefreshPostRequestBody()

	pausePeriod := int32(device.PauseTimePeriodInMinutes.ValueInt64())
	requestBody.SetPauseTimePeriodInMinutes(&pausePeriod)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final managed device pause config refresh request for device %s", device.DeviceID.ValueString()), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDevicePauseConfig) *devicemanagement.ComanagedDevicesItemPauseConfigurationRefreshPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemPauseConfigurationRefreshPostRequestBody()

	pausePeriod := int32(device.PauseTimePeriodInMinutes.ValueInt64())
	requestBody.SetPauseTimePeriodInMinutes(&pausePeriod)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final co-managed device pause config refresh request for device %s", device.DeviceID.ValueString()), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
