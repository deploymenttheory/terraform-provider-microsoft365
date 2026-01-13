package graphBetaRevokeAppleVppLicenses

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// VppLicenseValidationResult contains the results of device validation
type VppLicenseValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	NonAppleManagedDevices      []string
	NonAppleComanagedDevices    []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDeviceIDs []string, comanagedDeviceIDs []string) (*VppLicenseValidationResult, error) {
	result := &VppLicenseValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		NonAppleManagedDevices:      make([]string, 0),
		NonAppleComanagedDevices:    make([]string, 0),
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

		// Check that device is iOS or iPadOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if os != "ios" && os != "ipados" {
				result.NonAppleManagedDevices = append(result.NonAppleManagedDevices,
					fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Managed device is not iOS/iPadOS", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
				})
				continue
			}
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

		// Check that device is iOS or iPadOS
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if os != "ios" && os != "ipados" {
				result.NonAppleComanagedDevices = append(result.NonAppleComanagedDevices,
					fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Co-managed device is not iOS/iPadOS", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
				})
				continue
			}
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
