package graphBetaRotateBitLockerKeys

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

func (a *RotateBitLockerKeysAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RotateBitLockerKeysActionModel

	tflog.Debug(ctx, "Starting BitLocker key rotation action", map[string]any{"action": ActionName})

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

	// Convert framework lists to Go slices
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
	tflog.Debug(ctx, "Processing devices for BitLocker key rotation", map[string]any{
		"managed_devices":   len(managedDeviceIDs),
		"comanaged_devices": len(comanagedDeviceIDs),
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
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Error(validationResult.NonWindowsManagedDevices, "managed device", "are not Windows devices. BitLocker key rotation only works on Windows devices").
			Error(validationResult.NonWindowsComanagedDevices, "co-managed device", "are not Windows devices. BitLocker key rotation only works on Windows devices")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("BitLocker key rotation", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)))

	// Process managed devices sequentially
	for _, deviceID := range managedDeviceIDs {
		err := a.rotateManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to rotate BitLocker keys on managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("BitLocker keys rotated")
			tflog.Info(ctx, "Successfully rotated BitLocker keys on managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, deviceID := range comanagedDeviceIDs {
		err := a.rotateComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to rotate BitLocker keys on co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("BitLocker keys rotated")
			tflog.Info(ctx, "Successfully rotated BitLocker keys on co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("BitLocker key rotation")
			tflog.Warn(ctx, "BitLocker key rotation action completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("BitLocker Key Rotation Failed", "rotate BitLocker keys on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("rotated all BitLocker recovery keys. New recovery passwords have been generated and escrowed to Intune/Azure AD. Previous BitLocker recovery keys are now invalid. New recovery keys are available in the Microsoft Intune admin center and Azure AD portal")
	}

	tflog.Info(ctx, "BitLocker key rotation action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *RotateBitLockerKeysAction) rotateManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Rotating BitLocker keys on managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RotateBitLockerKeys().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *RotateBitLockerKeysAction) rotateComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Rotating BitLocker keys on co-managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RotateBitLockerKeys().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
