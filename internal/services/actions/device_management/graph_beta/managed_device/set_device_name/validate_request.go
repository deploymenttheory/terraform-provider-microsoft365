package graphBetaSetDeviceNameManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// SetDeviceNameValidationResult contains the results of device validation
type SetDeviceNameValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceSetName, comanagedDevices []ComanagedDeviceSetName) (*SetDeviceNameValidationResult, error) {
	result := &SetDeviceNameValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, device := range managedDevices {
		deviceID := device.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating managed device", map[string]any{"device_id": deviceID})

		_, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, deviceID)
				tflog.Warn(ctx, "Managed device not found", map[string]any{"device_id": deviceID})
				continue
			}
			return nil, fmt.Errorf("failed to validate managed device %s: %w", deviceID, err)
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices using managedDevices endpoint
	for _, device := range comanagedDevices {
		deviceID := device.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating co-managed device", map[string]any{"device_id": deviceID})

		_, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, deviceID)
				tflog.Warn(ctx, "Co-managed device not found", map[string]any{"device_id": deviceID})
				continue
			}
			return nil, fmt.Errorf("failed to validate co-managed device %s: %w", deviceID, err)
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}

