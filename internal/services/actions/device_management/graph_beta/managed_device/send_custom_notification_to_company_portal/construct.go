package graphBetaSendCustomNotificationToCompanyPortal

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceNotification) *devicemanagement.ManagedDevicesItemSendCustomNotificationToCompanyPortalPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemSendCustomNotificationToCompanyPortalPostRequestBody()

	convert.FrameworkToGraphString(device.NotificationTitle, requestBody.SetNotificationTitle)
	convert.FrameworkToGraphString(device.NotificationBody, requestBody.SetNotificationBody)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device notification request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceNotification) *devicemanagement.ComanagedDevicesItemSendCustomNotificationToCompanyPortalPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemSendCustomNotificationToCompanyPortalPostRequestBody()

	convert.FrameworkToGraphString(device.NotificationTitle, requestBody.SetNotificationTitle)
	convert.FrameworkToGraphString(device.NotificationBody, requestBody.SetNotificationBody)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device notification request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
