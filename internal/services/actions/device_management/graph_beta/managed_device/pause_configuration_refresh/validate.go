package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *PauseConfigurationRefreshManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data PauseConfigurationRefreshManagedDeviceActionModel

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

	tflog.Debug(ctx, fmt.Sprintf("Validating pause configuration refresh action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string

	// Validate managed devices
	for _, managedDevice := range data.ManagedDevices {
		deviceID := managedDevice.DeviceID.ValueString()

		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentManagedDevices = append(nonExistentManagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("managed_devices"),
					"Error Validating Managed Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility - configuration refresh pause only works on Windows devices
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
				}
			} else {
				unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}

			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully with pause period of %d minutes",
				deviceID, managedDevice.PauseTimePeriodInMinutes.ValueInt64()))
		}
	}

	// Validate co-managed devices
	for _, comanagedDevice := range data.ComanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()

		device, err := a.client.
			DeviceManagement().
			ComanagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentComanagedDevices = append(nonExistentComanagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("comanaged_devices"),
					"Error Validating Co-Managed Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
				}
			} else {
				unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}

			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully with pause period of %d minutes",
				deviceID, comanagedDevice.PauseTimePeriodInMinutes.ValueInt64()))
		}
	}

	if len(nonExistentManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Non-Existent Managed Devices",
			fmt.Sprintf("The following managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentManagedDevices, ", ")),
		)
	}

	if len(nonExistentComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Non-Existent Co-Managed Devices",
			fmt.Sprintf("The following co-managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing co-managed devices.",
				strings.Join(nonExistentComanagedDevices, ", ")),
		)
	}

	if len(unsupportedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Unsupported Managed Device OS",
			fmt.Sprintf("The following managed devices are not running a supported Windows operating system: %s. "+
				"Configuration refresh pause is only supported on Windows 10/11 devices.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Unsupported Co-Managed Device OS",
			fmt.Sprintf("The following co-managed devices are not running a supported Windows operating system: %s. "+
				"Configuration refresh pause is only supported on Windows 10/11 devices.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	// Add informational warning about the pause operation
	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	if totalDevices > 0 {
		resp.Diagnostics.AddWarning(
			"Configuration Refresh Pause",
			fmt.Sprintf("This action will pause configuration refresh for %d device(s).\n\n"+
				"Important notes:\n"+
				"- Devices will not receive new policy updates during the pause period\n"+
				"- Existing applied policies remain in effect\n"+
				"- Configuration refresh automatically resumes after the pause period\n"+
				"- Users can still manually sync from Company Portal\n"+
				"- Critical security updates may still be applied\n"+
				"- Use this feature during maintenance windows or troubleshooting only",
				totalDevices),
		)
	}
}

