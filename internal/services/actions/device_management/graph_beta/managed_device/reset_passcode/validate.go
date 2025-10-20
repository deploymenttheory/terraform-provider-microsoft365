package graphBetaResetManagedDevicePasscode

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ResetManagedDevicePasscodeAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ResetManagedDevicePasscodeActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that device_ids is not empty
	if data.DeviceIDs.IsNull() || data.DeviceIDs.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Missing Required Configuration",
			"device_ids must be specified and contain at least one device ID.",
		)
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(deviceIDs) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Invalid Configuration",
			"device_ids must contain at least one device ID.",
		)
		return
	}

	deviceIDMap := make(map[string]bool)
	var duplicates []string

	for _, deviceID := range deviceIDs {
		if deviceIDMap[deviceID] {
			duplicates = append(duplicates, deviceID)
		}
		deviceIDMap[deviceID] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs",
			fmt.Sprintf("The following device IDs are duplicated in the list: %v. "+
				"Each device passcode will only be reset once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating passcode reset action for %d device(s)", len(deviceIDs)))

	var unsupportedDevices []string

	for _, deviceID := range deviceIDs {
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("device_ids"),
				"Unable to Validate Device",
				fmt.Sprintf("Could not fetch managed device with ID %s to validate. "+
					"Ensure the device exists and you have permission to manage it. "+
					"Error: %s", deviceID, err.Error()),
			)
			continue
		}

		if device == nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("device_ids"),
				"Device Not Found",
				fmt.Sprintf("Managed device with ID %s was not found. "+
					"Ensure the device ID is correct and the device is enrolled in Intune.", deviceID),
			)
			continue
		}

		// Check platform compatibility - reset passcode is only supported on Android
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if os != "android" {
				unsupportedDevices = append(unsupportedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				continue
			}
		} else {
			unsupportedDevices = append(unsupportedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		// Log device details for debugging
		deviceName := "unknown"
		if device.GetDeviceName() != nil {
			deviceName = *device.GetDeviceName()
		}
		tflog.Debug(ctx, fmt.Sprintf("Validated device %s (Name: %s) for passcode reset", deviceID, deviceName))
	}

	// Error for unsupported devices
	if len(unsupportedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Unsupported Devices for Reset Passcode",
			fmt.Sprintf("Reset passcode is only supported on Android devices. The following devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedDevices, ", ")),
		)
	}

	// General warning about device connectivity
	if len(deviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"Device Connectivity Required",
			fmt.Sprintf("Resetting passcodes for %d device(s). "+
				"Devices must be online and connected to receive the reset passcode command. "+
				"If a device is offline, the command will be queued and executed when the device next checks in. "+
				"After successful reset, you must retrieve the new temporary passcodes from the Intune portal "+
				"and communicate them securely to the device users.",
				len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation complete for %d device(s)", len(deviceIDs)))
}
