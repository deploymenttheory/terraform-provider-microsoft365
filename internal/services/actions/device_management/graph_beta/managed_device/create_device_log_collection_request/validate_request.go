package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// LogCollectionValidationResult contains the results of device validation
type LogCollectionValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceLogCollection, comanagedDevices []ComanagedDeviceLogCollection) (*LogCollectionValidationResult, error) {
	result := &LogCollectionValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, managedDevice := range managedDevices {
		deviceID := managedDevice.DeviceID.ValueString()
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

		// Validate device supports log collection (Windows 10 1709+, Windows 11)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			// Log collection only supported on Windows devices
			isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
				deviceType == models.WINDOWSRT_DEVICETYPE ||
				deviceType == models.WINDOWS10X_DEVICETYPE ||
				deviceType == models.CLOUDPC_DEVICETYPE

			if !isWindowsDevice {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			} else {
				// Verify Windows 10 version 1709 or later
				if device.GetOsVersion() != nil {
					osVersion := *device.GetOsVersion()
					// Windows 10 version 1709 is build 10.0.16299
					// We check for version 10.0 and build >= 16299
					if strings.HasPrefix(osVersion, "10.0.") {
						versionParts := strings.Split(osVersion, ".")
						if len(versionParts) >= 3 {
							buildNumber := versionParts[2]
							// Simple string comparison works for build numbers
							if buildNumber < "16299" {
								result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
									fmt.Sprintf("%s (osVersion: %s, requires Windows 10 version 1709 or later)", deviceID, osVersion))
							}
						}
					}
				}
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices - use managedDevices endpoint as comanagedDevices/{id} doesn't support GET
	for _, comanagedDevice := range comanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()
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

		// Validate device supports log collection (Windows 10 1709+, Windows 11)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			// Log collection only supported on Windows devices
			isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
				deviceType == models.WINDOWSRT_DEVICETYPE ||
				deviceType == models.WINDOWS10X_DEVICETYPE ||
				deviceType == models.CLOUDPC_DEVICETYPE

			if !isWindowsDevice {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			} else {
				// Verify Windows 10 version 1709 or later
				if device.GetOsVersion() != nil {
					osVersion := *device.GetOsVersion()
					// Windows 10 version 1709 is build 10.0.16299
					if strings.HasPrefix(osVersion, "10.0.") {
						versionParts := strings.Split(osVersion, ".")
						if len(versionParts) >= 3 {
							buildNumber := versionParts[2]
							if buildNumber < "16299" {
								result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
									fmt.Sprintf("%s (osVersion: %s, requires Windows 10 version 1709 or later)", deviceID, osVersion))
							}
						}
					}
				}
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
