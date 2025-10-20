package graphBetaMoveDevicesToOUManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDevicesRequest(ctx context.Context, deviceIDs []string, ouPath string) *devicemanagement.ManagedDevicesMoveDevicesToOUPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesMoveDevicesToOUPostRequestBody()

	// Convert string device IDs to UUIDs
	deviceUUIDs := make([]uuid.UUID, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		if parsedUUID, err := uuid.Parse(deviceID); err == nil {
			deviceUUIDs = append(deviceUUIDs, parsedUUID)
		} else {
			tflog.Warn(ctx, "Failed to parse device ID as UUID", map[string]any{"deviceID": deviceID, "error": err.Error()})
		}
	}

	requestBody.SetDeviceIds(deviceUUIDs)
	requestBody.SetOrganizationalUnitPath(&ouPath)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed devices move to OU request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDevicesRequest(ctx context.Context, deviceIDs []string, ouPath string) *devicemanagement.ComanagedDevicesMoveDevicesToOUPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesMoveDevicesToOUPostRequestBody()

	// Convert string device IDs to UUIDs
	deviceUUIDs := make([]uuid.UUID, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		if parsedUUID, err := uuid.Parse(deviceID); err == nil {
			deviceUUIDs = append(deviceUUIDs, parsedUUID)
		} else {
			tflog.Warn(ctx, "Failed to parse device ID as UUID", map[string]any{"deviceID": deviceID, "error": err.Error()})
		}
	}

	requestBody.SetDeviceIds(deviceUUIDs)
	requestBody.SetOrganizationalUnitPath(&ouPath)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed devices move to OU request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
