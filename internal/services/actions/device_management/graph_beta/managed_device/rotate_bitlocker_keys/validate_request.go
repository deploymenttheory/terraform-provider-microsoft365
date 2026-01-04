package graphBetaRotateBitLockerKeys

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// BitLockerValidationResult contains the results of device validation
type BitLockerValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	NonWindowsManagedDevices    []string
	NonWindowsComanagedDevices  []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDeviceIDs []string, comanagedDeviceIDs []string) (*BitLockerValidationResult, error) {
	result := &BitLockerValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		NonWindowsManagedDevices:    make([]string, 0),
		NonWindowsComanagedDevices:  make([]string, 0),
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

		// Check that device is Windows
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices,
					fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Managed device is not Windows", map[string]any{
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

		// Check that device is Windows
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices,
					fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Co-managed device is not Windows", map[string]any{
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

