package graphBetaSendCustomNotificationToCompanyPortal

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// NotificationValidationResult contains the results of device validation
type NotificationValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceNotification, comanagedDevices []ComanagedDeviceNotification) (*NotificationValidationResult, error) {
	result := &NotificationValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, device := range managedDevices {
		deviceID := device.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating managed device", map[string]any{"device_id": deviceID})

		managedDevice, err := client.
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

		// Check platform compatibility - custom notifications only supported on iOS, iPadOS, and Android
		if managedDevice.GetOperatingSystem() != nil {
			os := strings.ToLower(*managedDevice.GetOperatingSystem())
			supportedOS := map[string]bool{
				"ios":     true,
				"ipados":  true,
				"android": true,
			}
			if !supportedOS[os] {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *managedDevice.GetOperatingSystem()))
				tflog.Warn(ctx, "Managed device OS not supported for custom notifications", map[string]any{
					"device_id": deviceID,
					"os":        *managedDevice.GetOperatingSystem(),
				})
				continue
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Managed device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices using managedDevices endpoint
	for _, device := range comanagedDevices {
		deviceID := device.DeviceID.ValueString()
		tflog.Debug(ctx, "Validating co-managed device", map[string]any{"device_id": deviceID})

		comanagedDevice, err := client.
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
		if comanagedDevice.GetOperatingSystem() != nil {
			os := strings.ToLower(*comanagedDevice.GetOperatingSystem())
			supportedOS := map[string]bool{
				"ios":     true,
				"ipados":  true,
				"android": true,
			}
			if !supportedOS[os] {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *comanagedDevice.GetOperatingSystem()))
				tflog.Warn(ctx, "Co-managed device OS not supported for custom notifications", map[string]any{
					"device_id": deviceID,
					"os":        *comanagedDevice.GetOperatingSystem(),
				})
				continue
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Co-managed device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}

