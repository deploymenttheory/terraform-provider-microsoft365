package graphBetaRetireManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// RetireValidationResult contains the results of device validation
type RetireValidationResult struct {
	NonExistentDevices []string
	UnsupportedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*RetireValidationResult, error) {
	result := &RetireValidationResult{
		NonExistentDevices: make([]string, 0),
		UnsupportedDevices: make([]string, 0),
	}

	for _, deviceID := range deviceIDs {
		tflog.Debug(ctx, "Validating device", map[string]any{"device_id": deviceID})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentDevices = append(result.NonExistentDevices, deviceID)
				tflog.Warn(ctx, "Device not found", map[string]any{"device_id": deviceID})
				continue
			}
			return nil, fmt.Errorf("failed to validate device %s: %w", deviceID, err)
		}

		// Check platform compatibility - retire is supported on Windows, iOS, iPadOS, macOS, Android, and ChromeOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			supportedOS := map[string]bool{
				"windows":  true,
				"ios":      true,
				"ipados":   true,
				"macos":    true,
				"android":  true,
				"chromeos": true,
			}
			if !supportedOS[os] {
				result.UnsupportedDevices = append(result.UnsupportedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Device has unsupported OS for retire", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
				})
				continue
			}
		} else {
			result.UnsupportedDevices = append(result.UnsupportedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
