package graphBetaPauseConfigurationRefreshManagedDevice

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

func (a *PauseConfigurationRefreshManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data PauseConfigurationRefreshManagedDeviceActionModel

	tflog.Debug(ctx, "Starting pause configuration refresh action", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for configuration refresh pause", map[string]any{
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not Windows devices (configuration refresh pause only works on Windows 10/11)").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Windows devices (configuration refresh pause only works on Windows 10/11)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("configuration refresh pause", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices sequentially
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		pauseMinutes := device.PauseTimePeriodInMinutes.ValueInt64()

		err := a.pauseConfigManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to pause configuration refresh for managed device", map[string]any{
				"device_id":     deviceID,
				"pause_minutes": pauseMinutes,
				"error":         err.Error(),
			})
		} else {
			// Format time description
			timeDesc := formatPauseTime(pauseMinutes)
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("configuration refresh paused for %s", timeDesc))
			tflog.Info(ctx, "Successfully paused configuration refresh for managed device", map[string]any{
				"device_id":     deviceID,
				"pause_minutes": pauseMinutes,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		pauseMinutes := device.PauseTimePeriodInMinutes.ValueInt64()

		err := a.pauseConfigComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to pause configuration refresh for co-managed device", map[string]any{
				"device_id":     deviceID,
				"pause_minutes": pauseMinutes,
				"error":         err.Error(),
			})
		} else {
			// Format time description
			timeDesc := formatPauseTime(pauseMinutes)
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("configuration refresh paused for %s", timeDesc))
			tflog.Info(ctx, "Successfully paused configuration refresh for co-managed device", map[string]any{
				"device_id":     deviceID,
				"pause_minutes": pauseMinutes,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("configuration refresh pause")
			tflog.Warn(ctx, "Configuration refresh pause completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Configuration Refresh Pause Failed", "pause configuration refresh for devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("paused configuration refresh for all devices. Configuration refresh will automatically resume after the specified time periods")
	}

	tflog.Info(ctx, "Pause configuration refresh action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *PauseConfigurationRefreshManagedDeviceAction) pauseConfigManagedDevice(ctx context.Context, device ManagedDevicePauseConfig) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Pausing configuration refresh for managed device", map[string]any{
		"device_id":     deviceID,
		"pause_minutes": device.PauseTimePeriodInMinutes.ValueInt64(),
	})

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		PauseConfigurationRefresh().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *PauseConfigurationRefreshManagedDeviceAction) pauseConfigComanagedDevice(ctx context.Context, device ComanagedDevicePauseConfig) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Pausing configuration refresh for co-managed device", map[string]any{
		"device_id":     deviceID,
		"pause_minutes": device.PauseTimePeriodInMinutes.ValueInt64(),
	})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		PauseConfigurationRefresh().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

// formatPauseTime formats pause duration for user-friendly display
func formatPauseTime(minutes int64) string {
	hours := minutes / 60
	mins := minutes % 60

	if hours > 0 && mins > 0 {
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d hour(s)", hours)
	}
	return fmt.Sprintf("%d minute(s)", minutes)
}
