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

	// Validate that at least one device list is provided
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device entry.",
		)
		return
	}

	// Check for duplicate device IDs in managed devices
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

	// Check for duplicate device IDs in co-managed devices
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

	// Check for devices appearing in both lists
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

	// Validate notification content for managed devices
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

	// Validate notification content for co-managed devices
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

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(data.ManagedDevices),
		"comanaged_count": len(data.ComanagedDevices),
		"total_devices":   totalDevices,
	})
}
