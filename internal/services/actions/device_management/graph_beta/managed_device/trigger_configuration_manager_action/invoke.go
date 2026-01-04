package graphBetaTriggerConfigurationManagerActionManagedDevice

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

func (a *TriggerConfigurationManagerActionManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TriggerConfigurationManagerActionManagedDeviceActionModel

	tflog.Debug(ctx, "Starting Configuration Manager action trigger", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for Configuration Manager actions", map[string]any{
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not supported for Configuration Manager actions. Configuration Manager actions require Windows devices with the Configuration Manager client installed").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not supported for Configuration Manager actions. Configuration Manager actions require Windows devices with the Configuration Manager client installed")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("Configuration Manager action triggers", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices sequentially
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		action := device.Action.ValueString()

		err := a.triggerConfigManagerActionManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("Configuration Manager action '%s' trigger failed: %s", action, err.Error()))
			tflog.Error(ctx, "Failed to trigger Configuration Manager action on managed device", map[string]any{
				"device_id": deviceID,
				"action":    action,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("Configuration Manager action '%s' triggered", action))
			tflog.Info(ctx, "Successfully triggered Configuration Manager action on managed device", map[string]any{
				"device_id": deviceID,
				"action":    action,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		action := device.Action.ValueString()

		err := a.triggerConfigManagerActionComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("Configuration Manager action '%s' trigger failed: %s", action, err.Error()))
			tflog.Error(ctx, "Failed to trigger Configuration Manager action on co-managed device", map[string]any{
				"device_id": deviceID,
				"action":    action,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("Configuration Manager action '%s' triggered", action))
			tflog.Info(ctx, "Successfully triggered Configuration Manager action on co-managed device", map[string]any{
				"device_id": deviceID,
				"action":    action,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("Configuration Manager action triggers")
			tflog.Warn(ctx, "Configuration Manager action trigger completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Configuration Manager Action Trigger Failed", "trigger Configuration Manager actions on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("triggered Configuration Manager actions on all devices. Devices will execute the requested actions when they receive the command (requires devices to be online)")
	}

	tflog.Info(ctx, "Configuration Manager action trigger completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) triggerConfigManagerActionManagedDevice(ctx context.Context, device ManagedDeviceConfigManagerAction) error {
	deviceID := device.DeviceID.ValueString()
	action := device.Action.ValueString()

	tflog.Debug(ctx, "Triggering Configuration Manager action on managed device", map[string]any{
		"device_id": deviceID,
		"action":    action,
	})

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		TriggerConfigurationManagerAction().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) triggerConfigManagerActionComanagedDevice(ctx context.Context, device ComanagedDeviceConfigManagerAction) error {
	deviceID := device.DeviceID.ValueString()
	action := device.Action.ValueString()

	tflog.Debug(ctx, "Triggering Configuration Manager action on co-managed device", map[string]any{
		"device_id": deviceID,
		"action":    action,
	})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		TriggerConfigurationManagerAction().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
