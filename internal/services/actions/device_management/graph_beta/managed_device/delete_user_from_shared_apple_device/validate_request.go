package graphBetaDeleteUserFromSharedAppleDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// SharedAppleDeviceValidationResult contains the results of device validation
type SharedAppleDeviceValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceUserPair, comanagedDevices []ComanagedDeviceUserPair) (*SharedAppleDeviceValidationResult, error) {
	result := &SharedAppleDeviceValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, deviceUser := range managedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating managed device", map[string]any{"device_id": deviceID})

		device, err := client.
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

		// Validate device is a Shared iPad (iOS/iPadOS with shared device mode)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			// Shared iPad mode is only supported on iOS/iPadOS devices
			if deviceType != models.IPAD_DEVICETYPE {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices - use managedDevices endpoint as comanagedDevices/{id} doesn't support GET
	for _, deviceUser := range comanagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating co-managed device", map[string]any{"device_id": deviceID})

		device, err := client.
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

		// Validate device is a Shared iPad (iOS/iPadOS with shared device mode)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			// Shared iPad mode is only supported on iOS/iPadOS devices
			if deviceType != models.IPAD_DEVICETYPE {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
