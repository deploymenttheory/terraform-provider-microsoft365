package graphBetaPlayLostModeSoundManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// PlayLostModeSoundValidationResult contains the results of device validation
type PlayLostModeSoundValidationResult struct {
	NonExistentManagedDevices      []string
	NonExistentComanagedDevices    []string
	UnsupportedManagedDevices      []string
	UnsupportedComanagedDevices    []string
	UnsupervisedManagedDevices     []string
	UnsupervisedComanagedDevices   []string
	NotInLostModeManagedDevices    []string
	NotInLostModeComanagedDevices  []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDevicePlaySound, comanagedDevices []ComanagedDevicePlaySound) (*PlayLostModeSoundValidationResult, error) {
	result := &PlayLostModeSoundValidationResult{
		NonExistentManagedDevices:      make([]string, 0),
		NonExistentComanagedDevices:    make([]string, 0),
		UnsupportedManagedDevices:      make([]string, 0),
		UnsupportedComanagedDevices:    make([]string, 0),
		UnsupervisedManagedDevices:     make([]string, 0),
		UnsupervisedComanagedDevices:   make([]string, 0),
		NotInLostModeManagedDevices:    make([]string, 0),
		NotInLostModeComanagedDevices:  make([]string, 0),
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

		// Check platform compatibility - lost mode sound only works on iOS/iPadOS devices
		if device.GetDeviceType() != nil {
			deviceType := device.GetDeviceType().String()
			if deviceType != "iPad" && deviceType != "iPhone" && deviceType != "iPod" {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Type: %s)", deviceID, deviceType))
				continue
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown device type)", deviceID))
			continue
		}

		// Check if device is supervised (required for lost mode)
		if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
			result.UnsupervisedManagedDevices = append(result.UnsupervisedManagedDevices, deviceID)
		}

		// Check if device is in lost mode
		if device.GetLostModeState() != nil {
			lostModeState := device.GetLostModeState().String()
			if lostModeState == "disabled" {
				result.NotInLostModeManagedDevices = append(result.NotInLostModeManagedDevices, fmt.Sprintf("%s (state: %s)", deviceID, lostModeState))
			}
		}

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices using managedDevices endpoint
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

		// Check platform compatibility - lost mode sound only works on iOS/iPadOS devices
		if device.GetDeviceType() != nil {
			deviceType := device.GetDeviceType().String()
			if deviceType != "iPad" && deviceType != "iPhone" && deviceType != "iPod" {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Type: %s)", deviceID, deviceType))
				continue
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown device type)", deviceID))
			continue
		}

		// Check if device is supervised (required for lost mode)
		if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
			result.UnsupervisedComanagedDevices = append(result.UnsupervisedComanagedDevices, deviceID)
		}

		// Check if device is in lost mode
		if device.GetLostModeState() != nil {
			lostModeState := device.GetLostModeState().String()
			if lostModeState == "disabled" {
				result.NotInLostModeComanagedDevices = append(result.NotInLostModeComanagedDevices, fmt.Sprintf("%s (state: %s)", deviceID, lostModeState))
			}
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}

