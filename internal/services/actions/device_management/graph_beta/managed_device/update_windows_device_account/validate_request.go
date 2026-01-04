package graphBetaUpdateWindowsDeviceAccount

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// DeviceAccountValidationResult contains the results of device validation
type DeviceAccountValidationResult struct {
	NonExistentManagedDevices   []string
	NonWindowsManagedDevices    []string
	NonExistentComanagedDevices []string
	NonWindowsComanagedDevices  []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceAccount, comanagedDevices []ComanagedDeviceAccount) (*DeviceAccountValidationResult, error) {
	result := &DeviceAccountValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonWindowsManagedDevices:    make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		NonWindowsComanagedDevices:  make([]string, 0),
	}

	// Validate managed devices
	for _, deviceConfig := range managedDevices {
		deviceID := deviceConfig.DeviceID.ValueString()
		email := deviceConfig.DeviceAccountEmail.ValueString()

		tflog.Debug(ctx, "Validating managed device", map[string]any{
			"device_id": deviceID,
			"email":     email,
		})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, fmt.Sprintf("%s (account: %s)", deviceID, email))
				tflog.Warn(ctx, "Managed device not found", map[string]any{
					"device_id": deviceID,
					"email":     email,
				})
				continue
			}
			return nil, fmt.Errorf("failed to validate managed device %s: %w", deviceID, err)
		}

		// Check that device is Windows
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (OS: %s, account: %s)", deviceID, *device.GetOperatingSystem(), email))
				tflog.Warn(ctx, "Managed device is not Windows", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
					"email":     email,
				})
				continue
			}
		} else {
			result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (Unknown OS, account: %s)", deviceID, email))
			tflog.Warn(ctx, "Managed device has unknown OS", map[string]any{
				"device_id": deviceID,
				"email":     email,
			})
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{
			"device_id": deviceID,
			"email":     email,
		})
	}

	// Validate co-managed devices using managedDevices endpoint
	for _, deviceConfig := range comanagedDevices {
		deviceID := deviceConfig.DeviceID.ValueString()
		email := deviceConfig.DeviceAccountEmail.ValueString()

		tflog.Debug(ctx, "Validating co-managed device", map[string]any{
			"device_id": deviceID,
			"email":     email,
		})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, fmt.Sprintf("%s (account: %s)", deviceID, email))
				tflog.Warn(ctx, "Co-managed device not found", map[string]any{
					"device_id": deviceID,
					"email":     email,
				})
				continue
			}
			return nil, fmt.Errorf("failed to validate co-managed device %s: %w", deviceID, err)
		}

		// Check that device is Windows
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (OS: %s, account: %s)", deviceID, *device.GetOperatingSystem(), email))
				tflog.Warn(ctx, "Co-managed device is not Windows", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
					"email":     email,
				})
				continue
			}
		} else {
			result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (Unknown OS, account: %s)", deviceID, email))
			tflog.Warn(ctx, "Co-managed device has unknown OS", map[string]any{
				"device_id": deviceID,
				"email":     email,
			})
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{
			"device_id": deviceID,
			"email":     email,
		})
	}

	return result, nil
}

