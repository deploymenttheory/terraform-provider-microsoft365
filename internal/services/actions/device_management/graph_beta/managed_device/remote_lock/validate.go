package graphBetaRemoteLockManagedDevice

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RemoteLockManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RemoteLockManagedDeviceActionModel

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
				"Each device will only receive the lock command once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating remote lock action for %d device(s)", len(deviceIDs)))

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

		// Log device details for debugging
		deviceName := "unknown"
		if device.GetDeviceName() != nil {
			deviceName = *device.GetDeviceName()
		}

		operatingSystem := "unknown"
		if device.GetOperatingSystem() != nil {
			operatingSystem = *device.GetOperatingSystem()
		}

		tflog.Debug(ctx, fmt.Sprintf("Validated device %s (Name: %s, OS: %s) for remote lock", deviceID, deviceName, operatingSystem))
	}

	// General warning about device action
	if len(deviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"Device Lock Warning",
			fmt.Sprintf("You are about to remotely lock %d device(s). "+
				"Devices will lock immediately when they receive the command. "+
				"Users will need to enter their passcode to unlock. "+
				"Ensure you have authorization to lock these devices. "+
				"If devices are offline, the lock command will be queued and executed when they next check in with Intune.",
				len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation complete for %d device(s)", len(deviceIDs)))
}
