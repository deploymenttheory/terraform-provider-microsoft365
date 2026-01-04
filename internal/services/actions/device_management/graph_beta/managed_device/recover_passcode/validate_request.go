package graphBetaRecoverManagedDevicePasscode

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// RecoverPasscodeValidationResult contains the results of device validation
type RecoverPasscodeValidationResult struct {
	NonExistentDevices   []string
	UnsupportedOSDevices []string
	UnsupervisedDevices  []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*RecoverPasscodeValidationResult, error) {
	result := &RecoverPasscodeValidationResult{
		NonExistentDevices:   make([]string, 0),
		UnsupportedOSDevices: make([]string, 0),
		UnsupervisedDevices:  make([]string, 0),
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

		// Check OS compatibility - Passcode recovery primarily works on iOS/iPadOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if os != "ios" && os != "ipados" {
				result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
			} else {
				// For iOS/iPadOS, check if supervised (better passcode escrow support)
				if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
					result.UnsupervisedDevices = append(result.UnsupervisedDevices, deviceID)
				}
			}
		} else {
			result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}

