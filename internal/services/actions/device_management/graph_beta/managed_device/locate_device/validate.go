package graphBetaLocateManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *LocateManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data LocateManagedDeviceActionModel

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
				"Locate request will only be sent once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating locate device action for %d device(s)", len(deviceIDs)))

	var nonExistentDevices []string
	var offlineDevices []string
	var locationDisabledDevices []string
	var unsupportedOSDevices []string

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
			// Check if device is online
			// Note: Device can still receive the command even if offline; it will process when it comes online
			if device.GetDeviceRegistrationState() != nil {
				regState := device.GetDeviceRegistrationState().String()
				if regState == "notRegisteredPendingEnrollment" || regState == "notRegistered" {
					offlineDevices = append(offlineDevices, fmt.Sprintf("%s (state: %s)", deviceID, regState))
				}
			}

			// Check OS support for locate device
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				// Locate device is supported on Windows, iOS, iPadOS, and Android only
				supportedOS := []string{"windows", "ios", "ipados", "android"}
				isSupported := false
				for _, supported := range supportedOS {
					if strings.Contains(os, supported) {
						isSupported = true
						break
					}
				}
				if !isSupported {
					unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				}

				// For iOS/iPadOS, check if supervised (better location support)
				if (os == "ios" || os == "ipados") && (device.GetIsSupervised() == nil || !*device.GetIsSupervised()) {
					locationDisabledDevices = append(locationDisabledDevices, fmt.Sprintf("%s (iOS/iPadOS - not supervised)", deviceID))
				}
			} else {
				unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
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

	if len(offlineDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Offline or Unregistered Devices",
			fmt.Sprintf("The following devices may be offline or not properly registered: %s. "+
				"The locate request will be queued and processed when the device comes online and checks in with Intune.",
				strings.Join(offlineDevices, ", ")),
		)
	}

	if len(locationDisabledDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Devices with Limited Location Support",
			fmt.Sprintf("The following devices may have limited location capabilities: %s. "+
				"For iOS/iPadOS devices, supervision provides more reliable location reporting. "+
				"Ensure location services are enabled on the device.",
				strings.Join(locationDisabledDevices, ", ")),
		)
	}

	if len(unsupportedOSDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Devices with Unsupported or Unknown Operating Systems",
			fmt.Sprintf("The following devices have operating systems that do not support locate device: %s. "+
				"Location reporting is only supported on Windows, iOS, iPadOS, and Android devices. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedOSDevices, ", ")),
		)
	}
}
