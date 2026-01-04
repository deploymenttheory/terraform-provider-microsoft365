package graphBetaRebootNowManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// RebootNowValidationResult contains the results of device validation
type RebootNowValidationResult struct {
	NonExistentDevices   []string
	UnsupportedOSDevices []string
	OfflineDevices       []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*RebootNowValidationResult, error) {
	result := &RebootNowValidationResult{
		NonExistentDevices:   make([]string, 0),
		UnsupportedOSDevices: make([]string, 0),
		OfflineDevices:       make([]string, 0),
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

		// Check OS compatibility - Reboot is supported on Windows, macOS, ChromeOS, and supervised iOS/iPadOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			supportedOS := map[string]bool{
				"windows":  true,
				"macos":    true,
				"chromeos": true,
				"ios":      true,
				"ipados":   true,
			}
			if !supportedOS[os] {
				result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (OS: %s - reboot not supported)", deviceID, *device.GetOperatingSystem()))
			}
		} else {
			result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
		}

		// Warn if device is offline
		if device.GetDeviceRegistrationState() != nil {
			regState := device.GetDeviceRegistrationState().String()
			if regState == "notRegisteredPendingEnrollment" || regState == "notRegistered" {
				result.OfflineDevices = append(result.OfflineDevices, fmt.Sprintf("%s (state: %s)", deviceID, regState))
			}
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
