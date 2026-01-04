package graphBetaCleanWindowsManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// WindowsDeviceValidationResult contains the results of device validation
type WindowsDeviceValidationResult struct {
	NonExistentManagedDevices          []string
	NonExistentComanagedDevices        []string
	NonWindowsManagedDevices           []string
	NonWindowsComanagedDevices         []string
	UnsupportedVersionManagedDevices   []string
	UnsupportedVersionComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceCleanWindows, comanagedDevices []ComanagedDeviceCleanWindows) (*WindowsDeviceValidationResult, error) {
	result := &WindowsDeviceValidationResult{
		NonExistentManagedDevices:          make([]string, 0),
		NonExistentComanagedDevices:        make([]string, 0),
		NonWindowsManagedDevices:           make([]string, 0),
		NonWindowsComanagedDevices:         make([]string, 0),
		UnsupportedVersionManagedDevices:   make([]string, 0),
		UnsupportedVersionComanagedDevices: make([]string, 0),
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

		// Validate device type supports Windows clean operation
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
				deviceType == models.WINDOWSRT_DEVICETYPE ||
				deviceType == models.WINDOWS10X_DEVICETYPE ||
				deviceType == models.CLOUDPC_DEVICETYPE

			if !isWindowsDevice {
				result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			} else {
				// For desktop and windowsRT, verify OS version is Windows 10+
				if (deviceType == models.DESKTOP_DEVICETYPE || deviceType == models.WINDOWSRT_DEVICETYPE) &&
					device.GetOsVersion() != nil {
					osVersion := *device.GetOsVersion()
					if !strings.HasPrefix(osVersion, "10.") {
						result.UnsupportedVersionManagedDevices = append(result.UnsupportedVersionManagedDevices,
							fmt.Sprintf("%s (deviceType: %s, osVersion: %s)", deviceID, deviceType.String(), osVersion))
					}
				}
			}
		} else {
			result.NonWindowsManagedDevices = append(result.NonWindowsManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
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

		// Validate device type supports Windows clean operation
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
				deviceType == models.WINDOWSRT_DEVICETYPE ||
				deviceType == models.WINDOWS10X_DEVICETYPE ||
				deviceType == models.CLOUDPC_DEVICETYPE

			if !isWindowsDevice {
				result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
			} else {
				// For desktop and windowsRT, verify OS version is Windows 10+
				if (deviceType == models.DESKTOP_DEVICETYPE || deviceType == models.WINDOWSRT_DEVICETYPE) &&
					device.GetOsVersion() != nil {
					osVersion := *device.GetOsVersion()
					if !strings.HasPrefix(osVersion, "10.") {
						result.UnsupportedVersionComanagedDevices = append(result.UnsupportedVersionComanagedDevices,
							fmt.Sprintf("%s (deviceType: %s, osVersion: %s)", deviceID, deviceType.String(), osVersion))
					}
				}
			}
		} else {
			result.NonWindowsComanagedDevices = append(result.NonWindowsComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
