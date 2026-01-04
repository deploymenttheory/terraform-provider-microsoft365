package graphBetaUpdateWindowsDeviceAccount

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *UpdateWindowsDeviceAccountAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data UpdateWindowsDeviceAccountActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one device list is provided
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	// Check for duplicate device IDs in managed devices
	if len(data.ManagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ManagedDevices {
			id := device.DeviceID.ValueString()
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"Each device will only be updated once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for duplicate device IDs in co-managed devices
	if len(data.ComanagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ComanagedDevices {
			id := device.DeviceID.ValueString()
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_devices"),
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"Each device will only be updated once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for devices appearing in both lists
	managedIDs := make(map[string]bool)
	for _, device := range data.ManagedDevices {
		managedIDs[device.DeviceID.ValueString()] = true
	}

	for _, device := range data.ComanagedDevices {
		id := device.DeviceID.ValueString()
		if managedIDs[id] {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Device ID in Both Lists",
				fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
					"A device should only be in one list. The update will be attempted for both endpoints, "+
					"but one may fail if the device is not actually of that type.",
					id),
			)
		}
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)

	// Informational message about device account update behavior
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Windows Device Account Update Information",
		fmt.Sprintf("This action will update device account configuration on %d Windows device(s). "+
			"This action is designed for shared Windows devices like Surface Hub and Microsoft Teams Rooms. "+
			"The updated configuration includes device account credentials, Exchange server settings, and synchronization options. "+
			"Devices may require a reboot for all changes to take effect. "+
			"Ensure the device account exists in Exchange/Microsoft 365 and has appropriate licenses and permissions.",
			totalDevices),
	)

	// Critical warning about passwords
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Password Security Warning",
		fmt.Sprintf("You are updating device account passwords for %d device(s). "+
			"Passwords are transmitted securely to devices and stored encrypted. "+
			"However, ensure you are following security best practices: "+
			"(1) Use strong, complex passwords "+
			"(2) Enable password rotation when possible "+
			"(3) Store passwords securely in Terraform using sensitive values "+
			"(4) Rotate passwords regularly "+
			"(5) Follow your organization's password policies. "+
			"Consider enabling automatic password rotation for enhanced security.",
			totalDevices),
	)

	// Warning about device reboot requirement
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Device Reboot May Be Required",
		fmt.Sprintf("After updating device account configuration on %d device(s), "+
			"the devices may need to be rebooted for changes to take full effect. "+
			"This affects device availability. Plan updates during maintenance windows when possible. "+
			"Users will not be able to use the devices during reboot. "+
			"Devices will automatically reconnect to Exchange and Teams/Skype for Business after restart.",
			totalDevices),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(data.ManagedDevices),
		"comanaged_count": len(data.ComanagedDevices),
		"total_devices":   totalDevices,
	})
}
