package graphBetaLocateManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// LocationValidationResult contains the results of device validation
type LocationValidationResult struct {
	NonExistentDevices      []string
	OfflineDevices          []string
	LocationDisabledDevices []string
	UnsupportedOSDevices    []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*LocationValidationResult, error) {
	result := &LocationValidationResult{
		NonExistentDevices:      make([]string, 0),
		OfflineDevices:          make([]string, 0),
		LocationDisabledDevices: make([]string, 0),
		UnsupportedOSDevices:    make([]string, 0),
	}

	// Validate devices
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

		// Check if device is online
		if device.GetDeviceRegistrationState() != nil {
			regState := device.GetDeviceRegistrationState().String()
			if regState == "notRegisteredPendingEnrollment" || regState == "notRegistered" {
				result.OfflineDevices = append(result.OfflineDevices, fmt.Sprintf("%s (state: %s)", deviceID, regState))
			}
		}

		// Check OS support for locate device
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			// Locate device is supported on Windows, iOS, iPadOS, and Android only
			supportedOS := []string{"windows", "ios", "ipados", "android"}
			isSupported := false
			for _, supported := range supportedOS {
				if strings.Contains(os, supported) {
					isSupported = true
					break
				}
			}
			if !isSupported {
				result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
			}

			// For iOS/iPadOS, check if supervised (better location support)
			if (os == "ios" || os == "ipados") && (device.GetIsSupervised() == nil || !*device.GetIsSupervised()) {
				result.LocationDisabledDevices = append(result.LocationDisabledDevices, fmt.Sprintf("%s (iOS/iPadOS - not supervised)", deviceID))
			}
		} else {
			result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
