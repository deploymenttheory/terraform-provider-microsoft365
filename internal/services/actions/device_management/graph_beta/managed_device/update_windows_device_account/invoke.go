package graphBetaUpdateWindowsDeviceAccount

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

func (a *UpdateWindowsDeviceAccountAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data UpdateWindowsDeviceAccountActionModel

	tflog.Debug(ctx, "Starting Windows device account update", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for Windows device account update", map[string]any{
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
			Error(validationResult.NonWindowsManagedDevices, "managed device", "are not Windows devices. This action only works on Windows devices").
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Error(validationResult.NonWindowsComanagedDevices, "co-managed device", "are not Windows devices. This action only works on Windows devices")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("Windows device account updates", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices sequentially
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		email := device.DeviceAccountEmail.ValueString()

		err := a.updateManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("device account update failed (account: %s): %s", email, err.Error()))
			tflog.Error(ctx, "Failed to update device account on managed device", map[string]any{
				"device_id": deviceID,
				"email":     email,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("device account updated (account: %s)", email))
			tflog.Info(ctx, "Successfully updated device account on managed device", map[string]any{
				"device_id": deviceID,
				"email":     email,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		email := device.DeviceAccountEmail.ValueString()

		err := a.updateComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("device account update failed (account: %s): %s", email, err.Error()))
			tflog.Error(ctx, "Failed to update device account on co-managed device", map[string]any{
				"device_id": deviceID,
				"email":     email,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("device account updated (account: %s)", email))
			tflog.Info(ctx, "Successfully updated device account on co-managed device", map[string]any{
				"device_id": deviceID,
				"email":     email,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("Windows device account updates")
			tflog.Warn(ctx, "Windows device account update completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Windows Device Account Update Failed", "update device accounts on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("updated device accounts on all devices. Device account credentials, Exchange server settings, and synchronization options have been updated. Devices may require a reboot for all changes to take effect. Verify functionality after devices restart and reconnect to Exchange and Teams/Skype for Business services")
	}

	tflog.Info(ctx, "Windows device account update completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *UpdateWindowsDeviceAccountAction) updateManagedDevice(ctx context.Context, device ManagedDeviceAccount) error {
	deviceID := device.DeviceID.ValueString()
	email := device.DeviceAccountEmail.ValueString()

	tflog.Debug(ctx, "Updating device account on managed device", map[string]any{
		"device_id": deviceID,
		"email":     email,
	})

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		UpdateWindowsDeviceAccount().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *UpdateWindowsDeviceAccountAction) updateComanagedDevice(ctx context.Context, device ComanagedDeviceAccount) error {
	deviceID := device.DeviceID.ValueString()
	email := device.DeviceAccountEmail.ValueString()

	tflog.Debug(ctx, "Updating device account on co-managed device", map[string]any{
		"device_id": deviceID,
		"email":     email,
	})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		UpdateWindowsDeviceAccount().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
