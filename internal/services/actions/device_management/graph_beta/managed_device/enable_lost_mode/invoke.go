package graphBetaEnableLostModeManagedDevice

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

func (a *EnableLostModeManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data EnableLostModeManagedDeviceActionModel

	tflog.Debug(ctx, "Starting enable lost mode action", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for lost mode enable", map[string]any{
		"total_devices": totalDevices,
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
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not managed by Intune").
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not iOS/iPadOS devices (lost mode is only supported on iOS/iPadOS)").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not iOS/iPadOS devices (lost mode is only supported on iOS/iPadOS)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("lost mode enable",
			fmt.Sprintf("%d managed, %d co-managed", len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.enableLostModeManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to enable lost mode for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Managed").Succeeded("lost mode enabled successfully")
			tflog.Info(ctx, "Successfully enabled lost mode for managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.enableLostModeComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to enable lost mode for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Co-managed").Succeeded("lost mode enabled successfully")
			tflog.Info(ctx, "Successfully enabled lost mode for co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("lost mode enable")
			tflog.Warn(ctx, "Lost mode enable completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Lost Mode Enable Failed", "enable lost mode on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("enabled lost mode on all devices. Devices are now locked with custom messages")
	}

	tflog.Info(ctx, "Enable lost mode action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *EnableLostModeManagedDeviceAction) enableLostModeManagedDevice(ctx context.Context, device ManagedDeviceLostMode) error {
	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(device.DeviceID.ValueString()).
		EnableLostMode().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *EnableLostModeManagedDeviceAction) enableLostModeComanagedDevice(ctx context.Context, device ComanagedDeviceLostMode) error {
	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(device.DeviceID.ValueString()).
		EnableLostMode().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
