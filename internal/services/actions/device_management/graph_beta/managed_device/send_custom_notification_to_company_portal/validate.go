package graphBetaSendCustomNotificationToCompanyPortal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *SendCustomNotificationToCompanyPortalAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data SendCustomNotificationToCompanyPortalActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device entry.",
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
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"Notifications will only be sent once per device, but you should remove duplicates from your configuration.",
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
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_devices"),
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"Notifications will only be sent once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	managedDeviceIDs := make(map[string]bool)
	for _, device := range data.ManagedDevices {
		managedDeviceIDs[device.DeviceID.ValueString()] = true
	}

	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		if managedDeviceIDs[deviceID] {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_devices"),
				"Device ID in Both Lists",
				fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
					"A device should only be in one list. The notification will be attempted for both endpoints, "+
					"but one may fail if the device is not actually of that type.",
					deviceID),
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating custom notification action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string

	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		managedDevice, err := a.client.
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
		} else if managedDevice != nil {
			// Check platform compatibility - custom notifications only supported on iOS, iPadOS, and Android
			if managedDevice.GetOperatingSystem() != nil {
				os := strings.ToLower(*managedDevice.GetOperatingSystem())
				supportedOS := map[string]bool{
					"ios":     true,
					"ipados":  true,
					"android": true,
				}
				if !supportedOS[os] {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *managedDevice.GetOperatingSystem()))
					continue
				}
			} else {
				unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		comanagedDevice, err := a.client.
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
		} else if comanagedDevice != nil {
			// Check platform compatibility - custom notifications only supported on iOS, iPadOS, and Android
			if comanagedDevice.GetOperatingSystem() != nil {
				os := strings.ToLower(*comanagedDevice.GetOperatingSystem())
				supportedOS := map[string]bool{
					"ios":     true,
					"ipados":  true,
					"android": true,
				}
				if !supportedOS[os] {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *comanagedDevice.GetOperatingSystem()))
					continue
				}
			} else {
				unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully", deviceID))
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
			"Unsupported Managed Devices for Custom Notifications",
			fmt.Sprintf("Custom notifications are only supported on iOS, iPadOS, and Android devices. The following managed devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Unsupported Co-Managed Devices for Custom Notifications",
			fmt.Sprintf("Custom notifications are only supported on iOS, iPadOS, and Android devices. The following co-managed devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	// Validate notification content
	for i, device := range data.ManagedDevices {
		if device.NotificationTitle.IsNull() || device.NotificationTitle.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("managed_devices").AtListIndex(i).AtName("notification_title"),
				"Empty Notification Title",
				fmt.Sprintf("Managed device %s has an empty notification title. All notifications must have a title.",
					device.DeviceID.ValueString()),
			)
		}

		if device.NotificationBody.IsNull() || device.NotificationBody.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("managed_devices").AtListIndex(i).AtName("notification_body"),
				"Empty Notification Body",
				fmt.Sprintf("Managed device %s has an empty notification body. All notifications must have a body message.",
					device.DeviceID.ValueString()),
			)
		}
	}

	for i, device := range data.ComanagedDevices {
		if device.NotificationTitle.IsNull() || device.NotificationTitle.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("comanaged_devices").AtListIndex(i).AtName("notification_title"),
				"Empty Notification Title",
				fmt.Sprintf("Co-managed device %s has an empty notification title. All notifications must have a title.",
					device.DeviceID.ValueString()),
			)
		}

		if device.NotificationBody.IsNull() || device.NotificationBody.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("comanaged_devices").AtListIndex(i).AtName("notification_body"),
				"Empty Notification Body",
				fmt.Sprintf("Co-managed device %s has an empty notification body. All notifications must have a body message.",
					device.DeviceID.ValueString()),
			)
		}
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)

	// Informational warnings
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Company Portal Notification Requirements",
		fmt.Sprintf("This action will send custom notifications to %d device(s) via the Company Portal app. "+
			"Important requirements: "+
			"(1) Company Portal app must be installed on the device "+
			"(2) User must be signed into Company Portal "+
			"(3) Device must be enrolled in Intune "+
			"(4) Device must have network connectivity. "+
			"If any of these requirements are not met, the notification delivery may fail silently. "+
			"Users will see the notification when they open or check the Company Portal app.",
			totalDevices),
	)

	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Notification Best Practices",
		fmt.Sprintf("You are sending custom notifications to %d device(s). Best practices: "+
			"(1) Keep notification titles concise (50-60 characters) "+
			"(2) Keep notification bodies clear and actionable (200-300 characters) "+
			"(3) Consider user time zones when sending time-sensitive notifications "+
			"(4) Use appropriate language and tone for your organization "+
			"(5) Include clear call-to-action or next steps "+
			"(6) Avoid sending excessive notifications that may be perceived as spam. "+
			"Notifications appear in the Company Portal app notifications section.",
			totalDevices),
	)
}
