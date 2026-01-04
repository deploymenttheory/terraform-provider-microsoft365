package graphBetaTriggerConfigurationManagerActionManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ConfigManagerValidationResult contains the results of device validation
type ConfigManagerValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceConfigManagerAction, comanagedDevices []ComanagedDeviceConfigManagerAction) (*ConfigManagerValidationResult, error) {
	result := &ConfigManagerValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, deviceConfig := range managedDevices {
		deviceID := deviceConfig.DeviceID.ValueString()
		action := deviceConfig.Action.ValueString()

		tflog.Debug(ctx, "Validating managed device", map[string]any{
			"device_id": deviceID,
			"action":    action,
		})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, fmt.Sprintf("%s (action: %s)", deviceID, action))
				tflog.Warn(ctx, "Managed device not found", map[string]any{
					"device_id": deviceID,
					"action":    action,
				})
				continue
			}
			return nil, fmt.Errorf("failed to validate managed device %s: %w", deviceID, err)
		}

		// Check platform compatibility - Configuration Manager is Windows-only
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (OS: %s, action: %s)", deviceID, *device.GetOperatingSystem(), action))
				tflog.Warn(ctx, "Managed device OS not supported for Configuration Manager actions", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
					"action":    action,
				})
				continue
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS, action: %s)", deviceID, action))
			tflog.Warn(ctx, "Managed device has unknown OS", map[string]any{
				"device_id": deviceID,
				"action":    action,
			})
			continue
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{
			"device_id": deviceID,
			"action":    action,
		})
	}

	// Validate co-managed devices using managedDevices endpoint
	for _, deviceConfig := range comanagedDevices {
		deviceID := deviceConfig.DeviceID.ValueString()
		action := deviceConfig.Action.ValueString()

		tflog.Debug(ctx, "Validating co-managed device", map[string]any{
			"device_id": deviceID,
			"action":    action,
		})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, fmt.Sprintf("%s (action: %s)", deviceID, action))
				tflog.Warn(ctx, "Co-managed device not found", map[string]any{
					"device_id": deviceID,
					"action":    action,
				})
				continue
			}
			return nil, fmt.Errorf("failed to validate co-managed device %s: %w", deviceID, err)
		}

		// Check platform compatibility
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s, action: %s)", deviceID, *device.GetOperatingSystem(), action))
				tflog.Warn(ctx, "Co-managed device OS not supported for Configuration Manager actions", map[string]any{
					"device_id": deviceID,
					"os":        *device.GetOperatingSystem(),
					"action":    action,
				})
				continue
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS, action: %s)", deviceID, action))
			tflog.Warn(ctx, "Co-managed device has unknown OS", map[string]any{
				"device_id": deviceID,
				"action":    action,
			})
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{
			"device_id": deviceID,
			"action":    action,
		})
	}

	return result, nil
}

