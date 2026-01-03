package graphBetaRebootNowManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RebootNowManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RebootNowManagedDeviceActionModel

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
				"Reboot command will only be sent once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating reboot now action for %d device(s)", len(deviceIDs)))

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
			// Check OS compatibility - Reboot is supported on Windows, macOS, ChromeOS, and supervised iOS/iPadOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				supportedOS := map[string]bool{
					"windows":  true,
					"macos":    true,
					"chromeos": true,
					"ios":      true,
					"ipados":   true,
				}
				if !supportedOS[os] {
					unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (OS: %s - reboot not supported)", deviceID, *device.GetOperatingSystem()))
				}
			} else {
				unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
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
			"Unsupported Devices for Reboot",
			fmt.Sprintf("The following devices do not support remote reboot: %s. "+
				"Remote reboot is supported on Windows, macOS, ChromeOS, and supervised iOS/iPadOS devices.",
				strings.Join(unsupportedOSDevices, ", ")),
		)
	}

	if len(offlineDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Offline or Unregistered Devices",
			fmt.Sprintf("The following devices may be offline or not properly registered: %s. "+
				"The reboot command will be queued and executed when the device comes online and checks in with Intune.",
				strings.Join(offlineDevices, ", ")),
		)
	}

	// General warning about user impact
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"User Impact Warning",
		fmt.Sprintf("Rebooting %d device(s) will immediately restart them when online. "+
			"Users may lose unsaved work and active sessions will be terminated. "+
			"Consider scheduling this action during maintenance windows or notifying users in advance.",
			len(deviceIDs)),
	)
}
