package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data CreateDeviceLogCollectionRequestManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	if len(data.ManagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ManagedDevices {
			deviceID := device.DeviceID.ValueString()
			if seen[deviceID] {
				duplicates = append(duplicates, deviceID)
			}
			seen[deviceID] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddWarning(
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"Log collection will only be requested once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	if len(data.ComanagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ComanagedDevices {
			deviceID := device.DeviceID.ValueString()
			if seen[deviceID] {
				duplicates = append(duplicates, deviceID)
			}
			seen[deviceID] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddWarning(
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"Log collection will only be requested once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	for _, managedDevice := range data.ManagedDevices {
		managedID := managedDevice.DeviceID.ValueString()
		for _, comanagedDevice := range data.ComanagedDevices {
			if managedID == comanagedDevice.DeviceID.ValueString() {
				resp.Diagnostics.AddWarning(
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
						"A device should only be in one list. Log collection will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
	})
}
