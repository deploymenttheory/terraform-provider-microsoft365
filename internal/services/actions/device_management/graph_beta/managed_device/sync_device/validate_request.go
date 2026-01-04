package graphBetaSyncManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// SyncValidationResult contains the results of device validation
type SyncValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDeviceIDs []string, comanagedDeviceIDs []string) (*SyncValidationResult, error) {
	result := &SyncValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, deviceID := range managedDeviceIDs {
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

		// Check platform compatibility - sync is supported on Windows, macOS, iOS, iPadOS, and Android
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			supportedOS := map[string]bool{
				"windows": true,
				"macos":   true,
				"ios":     true,
				"ipados":  true,
				"android": true,
			}
			if !supportedOS[os] {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Managed device OS not supported for sync", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
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
	for _, deviceID := range comanagedDeviceIDs {
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
			os := strings.ToLower(*device.GetOperatingSystem())
			supportedOS := map[string]bool{
				"windows": true,
				"macos":   true,
				"ios":     true,
				"ipados":  true,
				"android": true,
			}
			if !supportedOS[os] {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Co-managed device OS not supported for sync", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
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
