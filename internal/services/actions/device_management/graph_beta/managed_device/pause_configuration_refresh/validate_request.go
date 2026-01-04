package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// PauseConfigValidationResult contains the results of device validation
type PauseConfigValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDevicePauseConfig, comanagedDevices []ComanagedDevicePauseConfig) (*PauseConfigValidationResult, error) {
	result := &PauseConfigValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, managedDevice := range managedDevices {
		deviceID := managedDevice.DeviceID.ValueString()
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

		// Check platform compatibility - configuration refresh pause only works on Windows devices
		if device.GetOperatingSystem() != nil {
			osName := *device.GetOperatingSystem()
			if !strings.Contains(strings.ToLower(osName), "windows") {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
				continue
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{
			"device_id":     deviceID,
			"pause_minutes": managedDevice.PauseTimePeriodInMinutes.ValueInt64(),
		})
	}

	// Validate co-managed devices
	for _, comanagedDevice := range comanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()
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

		// Check platform compatibility
		if device.GetOperatingSystem() != nil {
			osName := *device.GetOperatingSystem()
			if !strings.Contains(strings.ToLower(osName), "windows") {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
				continue
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{
			"device_id":     deviceID,
			"pause_minutes": comanagedDevice.PauseTimePeriodInMinutes.ValueInt64(),
		})
	}

	return result, nil
}

