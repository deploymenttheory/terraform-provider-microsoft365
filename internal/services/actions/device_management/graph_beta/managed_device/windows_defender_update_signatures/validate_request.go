package graphBetaWindowsDefenderUpdateSignatures

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type SignatureUpdateValidationResult struct {
	NonExistentManagedDevices   []string
	NonWindowsManagedDevices    []string
	NonExistentComanagedDevices []string
	NonWindowsComanagedDevices  []string
}

func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDeviceIDs, comanagedDeviceIDs []string) (*SignatureUpdateValidationResult, error) {
	result := &SignatureUpdateValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonWindowsManagedDevices:    make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		NonWindowsComanagedDevices:  make([]string, 0),
	}

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

		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Managed device is not Windows", map[string]any{"device_id": deviceID, "os": *device.GetOperatingSystem()})
				continue
			}
		} else {
			result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Managed device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

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

		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				tflog.Warn(ctx, "Co-managed device is not Windows", map[string]any{"device_id": deviceID, "os": *device.GetOperatingSystem()})
				continue
			}
		} else {
			result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			tflog.Warn(ctx, "Co-managed device has unknown OS", map[string]any{"device_id": deviceID})
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}

