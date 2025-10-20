package graphBetaDeprovisionManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceDeprovision) *devicemanagement.ManagedDevicesItemDeprovisionPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemDeprovisionPostRequestBody()

	convert.FrameworkToGraphString(device.DeprovisionReason, requestBody.SetDeprovisionReason)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device deprovision request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceDeprovision) *devicemanagement.ComanagedDevicesItemDeprovisionPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemDeprovisionPostRequestBody()

	convert.FrameworkToGraphString(device.DeprovisionReason, requestBody.SetDeprovisionReason)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device deprovision request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
