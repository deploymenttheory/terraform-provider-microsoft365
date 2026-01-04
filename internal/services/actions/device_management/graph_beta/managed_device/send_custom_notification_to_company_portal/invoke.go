package graphBetaSendCustomNotificationToCompanyPortal

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/progress"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validation"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *SendCustomNotificationToCompanyPortalAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SendCustomNotificationToCompanyPortalActionModel

	tflog.Debug(ctx, "Starting custom notification delivery action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle timeout
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for custom notification delivery", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
		"total_devices":     totalDevices,
	})

	// Get ignore_partial_failures setting
	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	// Get validate_device_exists setting (default: true)
	validateDeviceExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateDeviceExists = data.ValidateDeviceExists.ValueBool()
	}

	// Perform API validation of devices if enabled
	if validateDeviceExists {
		tflog.Debug(ctx, "Performing device validation via API")

		validationResult, err := validateRequest(ctx, a.client, data.ManagedDevices, data.ComanagedDevices)
		if err != nil {
			tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError(
				"Device Validation Failed",
				fmt.Sprintf("Failed to validate devices: %s", err.Error()),
			)
			return
		}

		// Report validation results
		results := validation.NewResults().
			Error(validationResult.NonExistentManagedDevices, "managed device", "do not exist or are not managed by Intune").
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not supported for custom notifications. Only iOS, iPadOS, and Android devices support Company Portal notifications").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not supported for custom notifications. Only iOS, iPadOS, and Android devices support Company Portal notifications")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("custom notification delivery", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices sequentially
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		notificationTitle := device.NotificationTitle.ValueString()

		err := a.sendToManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("notification delivery failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to send notification to managed device", map[string]any{
				"device_id":          deviceID,
				"notification_title": notificationTitle,
				"error":              err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("notification sent: \"%s\"", notificationTitle))
			tflog.Info(ctx, "Successfully sent notification to managed device", map[string]any{
				"device_id":          deviceID,
				"notification_title": notificationTitle,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		notificationTitle := device.NotificationTitle.ValueString()

		err := a.sendToComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("notification delivery failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to send notification to co-managed device", map[string]any{
				"device_id":          deviceID,
				"notification_title": notificationTitle,
				"error":              err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("notification sent: \"%s\"", notificationTitle))
			tflog.Info(ctx, "Successfully sent notification to co-managed device", map[string]any{
				"device_id":          deviceID,
				"notification_title": notificationTitle,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("custom notification delivery")
			tflog.Warn(ctx, "Custom notification delivery action completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Custom Notification Delivery Failed", "send notifications to devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("sent all custom notifications. Notifications have been delivered to the Company Portal app on the target devices. " +
			"Users will see the notifications when they open or check the Company Portal app. Notifications appear in the Company Portal notifications section " +
			"and may also trigger push notifications depending on the user's notification settings and device configuration")
	}

	tflog.Info(ctx, "Custom notification delivery action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *SendCustomNotificationToCompanyPortalAction) sendToManagedDevice(ctx context.Context, device ManagedDeviceNotification) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Sending custom notification to managed device", map[string]any{
		"device_id":          deviceID,
		"notification_title": device.NotificationTitle.ValueString(),
	})

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SendCustomNotificationToCompanyPortal().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *SendCustomNotificationToCompanyPortalAction) sendToComanagedDevice(ctx context.Context, device ComanagedDeviceNotification) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Sending custom notification to co-managed device", map[string]any{
		"device_id":          deviceID,
		"notification_title": device.NotificationTitle.ValueString(),
	})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SendCustomNotificationToCompanyPortal().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
