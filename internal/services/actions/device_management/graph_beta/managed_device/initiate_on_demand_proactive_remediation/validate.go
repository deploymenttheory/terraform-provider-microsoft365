package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data InitiateOnDemandProactiveRemediationManagedDeviceActionModel

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

	// Check for duplicate device IDs within managed devices
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
				fmt.Sprintf("The following managed device IDs are duplicated in managed_devices: %s. "+
					"Proactive remediation will only be performed once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for duplicate device IDs within co-managed devices
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
				fmt.Sprintf("The following co-managed device IDs are duplicated in comanaged_devices: %s. "+
					"Proactive remediation will only be performed once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for devices appearing in both lists
	if len(data.ManagedDevices) > 0 && len(data.ComanagedDevices) > 0 {
		for _, managedDevice := range data.ManagedDevices {
			managedID := managedDevice.DeviceID.ValueString()
			for _, comanagedDevice := range data.ComanagedDevices {
				comanagedID := comanagedDevice.DeviceID.ValueString()
				if managedID == comanagedID {
					resp.Diagnostics.AddWarning(
						"Device ID in Both Lists",
						fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
							"A device should only be in one list. Proactive remediation will be attempted for both endpoints, "+
							"but one may fail if the device is not actually of that type.",
							managedID),
					)
				}
			}
		}
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
	})
}
