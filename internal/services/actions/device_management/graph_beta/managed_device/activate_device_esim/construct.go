package graphBetaActivateDeviceEsimManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceActivateEsim) *devicemanagement.ManagedDevicesItemActivateDeviceEsimPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemActivateDeviceEsimPostRequestBody()

	convert.FrameworkToGraphString(device.CarrierURL, requestBody.SetCarrierUrl)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device activate eSIM request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceActivateEsim) *devicemanagement.ComanagedDevicesItemActivateDeviceEsimPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemActivateDeviceEsimPostRequestBody()

	convert.FrameworkToGraphString(device.CarrierURL, requestBody.SetCarrierUrl)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device activate eSIM request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

