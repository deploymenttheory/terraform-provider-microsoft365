package graphBetaWipeManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WipeValidationResult struct {
	NonExistentDevices          []string
	UnsupportedDevices          []string
	ActivationLockWarningIDs    []string
	ActivationLockWarningOSList []string
}

func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string, macOsUnlockCodeProvided bool) (*WipeValidationResult, error) {
	result := &WipeValidationResult{
		NonExistentDevices:          make([]string, 0),
		UnsupportedDevices:          make([]string, 0),
		ActivationLockWarningIDs:    make([]string, 0),
		ActivationLockWarningOSList: make([]string, 0),
	}

	supportedOS := map[string]bool{
		"windows": true,
		"ios":     true,
		"ipados":  true,
		"macos":   true,
		"android": true,
	}

	for _, deviceID := range deviceIDs {
		tflog.Debug(ctx, "Validating device for wipe", map[string]any{"device_id": deviceID})

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

		// Check platform compatibility
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !supportedOS[os] {
				result.UnsupportedDevices = append(result.UnsupportedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Device OS not supported for wipe", map[string]any{"device_id": deviceID, "os": *device.GetOperatingSystem()})
				continue
			}
		} else {
			result.UnsupportedDevices = append(result.UnsupportedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		// Check for activation lock on iOS/macOS devices
		if device.GetOperatingSystem() != nil {
			os := *device.GetOperatingSystem()
			if (os == "iOS" || os == "iPadOS" || os == "macOS") && !macOsUnlockCodeProvided {
				activationLockBypassCode := device.GetActivationLockBypassCode()
				if activationLockBypassCode != nil && *activationLockBypassCode != "" {
					result.ActivationLockWarningIDs = append(result.ActivationLockWarningIDs, deviceID)
					result.ActivationLockWarningOSList = append(result.ActivationLockWarningOSList, os)
					tflog.Warn(ctx, "Device may have Activation Lock enabled", map[string]any{"device_id": deviceID, "os": os})
				}
			}
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
