package graphBetaLogoutSharedAppleDeviceActiveUser

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// LogoutValidationResult contains the results of device validation
type LogoutValidationResult struct {
	NonExistentDevices          []string
	NonIPadDevices              []string
	UnsupervisedDevices         []string
	PotentiallyNotSharedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*LogoutValidationResult, error) {
	result := &LogoutValidationResult{
		NonExistentDevices:          make([]string, 0),
		NonIPadDevices:              make([]string, 0),
		UnsupervisedDevices:         make([]string, 0),
		PotentiallyNotSharedDevices: make([]string, 0),
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

		// Check that device is iPadOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if os != "ipados" && os != "ios" {
				result.NonIPadDevices = append(result.NonIPadDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
			} else {
				// For iPadOS/iOS devices, check if supervised
				if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
					result.UnsupervisedDevices = append(result.UnsupervisedDevices, deviceID)
				}

				// Note: We cannot directly verify if a device is in Shared iPad mode from the API
				// The action will simply fail gracefully if the device is not in Shared iPad mode
				// We can check for iOS devices (iPhones don't support Shared mode)
				if os == "ios" {
					result.PotentiallyNotSharedDevices = append(result.PotentiallyNotSharedDevices,
						fmt.Sprintf("%s (iOS - likely iPhone, not Shared iPad)", deviceID))
				}
			}
		} else {
			result.NonIPadDevices = append(result.NonIPadDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
		}

		tflog.Debug(ctx, "Device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
