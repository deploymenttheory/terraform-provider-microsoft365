package graphBetaSetDeviceNameManagedDevice

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

func (a *SetDeviceNameManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SetDeviceNameManagedDeviceActionModel

	tflog.Debug(ctx, "Starting set device name action", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for set device name", map[string]any{
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
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("device name change", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices sequentially
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		deviceName := device.DeviceName.ValueString()

		err := a.setDeviceNameManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("device name change failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to set device name for managed device", map[string]any{
				"device_id":   deviceID,
				"device_name": deviceName,
				"error":       err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("device name changed to \"%s\"", deviceName))
			tflog.Info(ctx, "Successfully set device name for managed device", map[string]any{
				"device_id":   deviceID,
				"device_name": deviceName,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		deviceName := device.DeviceName.ValueString()

		err := a.setDeviceNameComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("device name change failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to set device name for co-managed device", map[string]any{
				"device_id":   deviceID,
				"device_name": deviceName,
				"error":       err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("device name changed to \"%s\"", deviceName))
			tflog.Info(ctx, "Successfully set device name for co-managed device", map[string]any{
				"device_id":   deviceID,
				"device_name": deviceName,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("device name change")
			tflog.Warn(ctx, "Set device name action completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Device Name Change Failed", "set device names")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("changed all device names. The name change command has been sent to all target devices. " +
			"Device names will be updated in the Intune console after devices check in and process the command. The time required for " +
			"name changes to reflect varies by platform and device online status")
	}

	tflog.Info(ctx, "Set device name action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *SetDeviceNameManagedDeviceAction) setDeviceNameManagedDevice(ctx context.Context, device ManagedDeviceSetName) error {
	deviceID := device.DeviceID.ValueString()
	deviceName := device.DeviceName.ValueString()
	tflog.Debug(ctx, "Setting device name for managed device", map[string]any{
		"device_id":   deviceID,
		"device_name": deviceName,
	})

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SetDeviceName().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *SetDeviceNameManagedDeviceAction) setDeviceNameComanagedDevice(ctx context.Context, device ComanagedDeviceSetName) error {
	deviceID := device.DeviceID.ValueString()
	deviceName := device.DeviceName.ValueString()
	tflog.Debug(ctx, "Setting device name for co-managed device", map[string]any{
		"device_id":   deviceID,
		"device_name": deviceName,
	})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SetDeviceName().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
