package graphBetaShutdownManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ShutdownManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ShutdownManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deviceIDs := make([]string, 0, len(data.DeviceIDs.Elements()))
	for _, elem := range data.DeviceIDs.Elements() {
		deviceIDs = append(deviceIDs, elem.String())
	}

	seen := make(map[string]bool)
	var duplicates []string
	for _, id := range deviceIDs {
		if seen[id] {
			duplicates = append(duplicates, id)
		}
		seen[id] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs Found",
			fmt.Sprintf("The following device IDs are duplicated in the configuration: %s. "+
				"Shutdown command will only be sent once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating shutdown action for %d device(s)", len(deviceIDs)))

	var nonExistentDevices []string
	var unsupportedOSDevices []string
	var offlineDevices []string

	for _, deviceID := range deviceIDs {
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentDevices = append(nonExistentDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("device_ids"),
					"Error Validating Device Existence",
					fmt.Sprintf("Failed to check existence of device %s: %s", deviceID, err.Error()),
				)
				return
			}
		} else {
			// Check OS compatibility - Shutdown is best supported on Windows and macOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				// Android doesn't support shutdown via Intune
				if strings.Contains(os, "android") {
					unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (OS: %s - shutdown not supported)", deviceID, *device.GetOperatingSystem()))
				}
			}

			// Warn if device is offline
			if device.GetDeviceRegistrationState() != nil {
				regState := device.GetDeviceRegistrationState().String()
				if regState == "notRegisteredPendingEnrollment" || regState == "notRegistered" {
					offlineDevices = append(offlineDevices, fmt.Sprintf("%s (state: %s)", deviceID, regState))
				}
			}
		}
	}

	if len(nonExistentDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Non-Existent Devices",
			fmt.Sprintf("The following device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentDevices, ", ")),
		)
	}

	if len(unsupportedOSDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Unsupported Devices for Shutdown",
			fmt.Sprintf("The following devices do not support remote shutdown: %s. "+
				"Remote shutdown is supported on Windows, macOS, and supervised iOS/iPadOS devices.",
				strings.Join(unsupportedOSDevices, ", ")),
		)
	}

	if len(offlineDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Offline or Unregistered Devices",
			fmt.Sprintf("The following devices may be offline or not properly registered: %s. "+
				"The shutdown command will be queued and executed when the device comes online and checks in with Intune.",
				strings.Join(offlineDevices, ", ")),
		)
	}

	// Critical warning about shutdown requiring manual power-on
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"Critical: Manual Power-On Required",
		fmt.Sprintf("Shutting down %d device(s) will POWER THEM OFF COMPLETELY. "+
			"Physical access will be required to power devices back on. "+
			"Users may lose unsaved work and will be unable to access their devices until manually powered on. "+
			"Consider using reboot action instead if devices need to come back online automatically. "+
			"Ensure you have legitimate business reason and proper authorization for this disruptive action.",
			len(deviceIDs)),
	)
}
