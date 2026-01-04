package graphBetaDisableLostModeManagedDevice

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

func (a *DisableLostModeManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DisableLostModeManagedDeviceActionModel

	tflog.Debug(ctx, "Starting disable lost mode action", map[string]any{"action": ActionName})

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

	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)
	tflog.Debug(ctx, "Processing devices for lost mode disable", map[string]any{
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

		validationResult, err := validateRequest(ctx, a.client, managedDeviceIDs, comanagedDeviceIDs)
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
		Starting("lost mode disable",
			fmt.Sprintf("%d managed, %d co-managed", len(managedDeviceIDs), len(comanagedDeviceIDs)))

	// Process managed devices
	for _, deviceID := range managedDeviceIDs {
		err := a.disableLostModeManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to disable lost mode for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Managed").Succeeded("lost mode disabled successfully")
			tflog.Info(ctx, "Successfully disabled lost mode for managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		err := a.disableLostModeComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to disable lost mode for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Co-managed").Succeeded("lost mode disabled successfully")
			tflog.Info(ctx, "Successfully disabled lost mode for co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("lost mode disable")
			tflog.Warn(ctx, "Lost mode disable completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Lost Mode Disable Failed", "disable lost mode on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("disabled lost mode on all devices. Devices have returned to normal operation")
	}

	tflog.Info(ctx, "Disable lost mode action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *DisableLostModeManagedDeviceAction) disableLostModeManagedDevice(ctx context.Context, deviceID string) error {
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		DisableLostMode().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *DisableLostModeManagedDeviceAction) disableLostModeComanagedDevice(ctx context.Context, deviceID string) error {
	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		DisableLostMode().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
