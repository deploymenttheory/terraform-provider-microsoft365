package graphBetaWindowsDefenderScan

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type DefenderScanValidationResult struct {
	NonExistentManagedDevices   []string
	NonWindowsManagedDevices    []string
	NonExistentComanagedDevices []string
	NonWindowsComanagedDevices  []string
}

func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceScan, comanagedDevices []ComanagedDeviceScan) (*DefenderScanValidationResult, error) {
	result := &DefenderScanValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonWindowsManagedDevices:    make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		NonWindowsComanagedDevices:  make([]string, 0),
	}

	for _, device := range managedDevices {
		deviceID := device.DeviceID.ValueString()
		scanType := "full"
		if device.QuickScan.ValueBool() {
			scanType = "quick"
		}

		tflog.Debug(ctx, "Validating managed device", map[string]any{
			"device_id": deviceID,
			"scan_type": scanType,
		})

		managedDevice, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, fmt.Sprintf("%s (%s scan)", deviceID, scanType))
				tflog.Warn(ctx, "Managed device not found", map[string]any{"device_id": deviceID, "scan_type": scanType})
				continue
			}
			return nil, fmt.Errorf("failed to validate managed device %s: %w", deviceID, err)
		}

		if managedDevice.GetOperatingSystem() != nil {
			os := strings.ToLower(*managedDevice.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (OS: %s, %s scan)", deviceID, *managedDevice.GetOperatingSystem(), scanType))
				tflog.Warn(ctx, "Managed device is not Windows", map[string]any{"device_id": deviceID, "os": *managedDevice.GetOperatingSystem(), "scan_type": scanType})
				continue
			}
		} else {
			result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (Unknown OS, %s scan)", deviceID, scanType))
			tflog.Warn(ctx, "Managed device has unknown OS", map[string]any{"device_id": deviceID, "scan_type": scanType})
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID, "scan_type": scanType})
	}

	for _, device := range comanagedDevices {
		deviceID := device.DeviceID.ValueString()
		scanType := "full"
		if device.QuickScan.ValueBool() {
			scanType = "quick"
		}

		tflog.Debug(ctx, "Validating co-managed device", map[string]any{
			"device_id": deviceID,
			"scan_type": scanType,
		})

		comanagedDevice, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, fmt.Sprintf("%s (%s scan)", deviceID, scanType))
				tflog.Warn(ctx, "Co-managed device not found", map[string]any{"device_id": deviceID, "scan_type": scanType})
				continue
			}
			return nil, fmt.Errorf("failed to validate co-managed device %s: %w", deviceID, err)
		}

		if comanagedDevice.GetOperatingSystem() != nil {
			os := strings.ToLower(*comanagedDevice.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (OS: %s, %s scan)", deviceID, *comanagedDevice.GetOperatingSystem(), scanType))
				tflog.Warn(ctx, "Co-managed device is not Windows", map[string]any{"device_id": deviceID, "os": *comanagedDevice.GetOperatingSystem(), "scan_type": scanType})
				continue
			}
		} else {
			result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (Unknown OS, %s scan)", deviceID, scanType))
			tflog.Warn(ctx, "Co-managed device has unknown OS", map[string]any{"device_id": deviceID, "scan_type": scanType})
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID, "scan_type": scanType})
	}

	return result, nil
}
