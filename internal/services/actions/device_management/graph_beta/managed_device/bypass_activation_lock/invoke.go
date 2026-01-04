package graphBetaBypassActivationLockManagedDevice

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

func (a *BypassActivationLockManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data BypassActivationLockManagedDeviceActionModel

	tflog.Debug(ctx, "Starting Activation Lock bypass action", map[string]any{"action": ActionName})

	// Parse configuration from the request
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, "Processing devices for Activation Lock bypass", map[string]any{
		"total_devices": totalDevices,
	})

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

		validationResult, err := constructBypassActivationLockRequest(ctx, a.client, deviceIDs)
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
			Error(validationResult.NonExistentDevices, "device", "do not exist or are not managed by Intune").
			Error(validationResult.UnsupportedOSDevices, "device", "have unsupported OS types (only iOS, iPadOS, and macOS are supported)").
			Error(validationResult.UnsupervisedIOSDevices, "iOS/iPadOS device", "are not supervised (supervision required)").
			Warning(validationResult.AlreadyBypassedDevices, "device", "already have bypass codes").
			Warning(validationResult.ActivationLockDisabledDevices, "device", "may not have Activation Lock enabled")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("Activation Lock bypass")

	// Process each device
	for _, deviceID := range deviceIDs {
		err := a.bypassActivationLock(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID).Failed(err.Error())
			tflog.Error(ctx, "Failed to bypass Activation Lock for device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID).Succeeded("Activation Lock bypass code generated successfully")
			tflog.Info(ctx, "Successfully bypassed Activation Lock for device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("Activation Lock bypass")
			tflog.Warn(ctx, "Activation Lock bypass completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Activation Lock Bypass Failed", "bypass Activation Lock")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("bypassed Activation Lock on all devices. Bypass codes are now available in Intune device properties")
	}

	tflog.Info(ctx, "Activation Lock bypass action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

// bypassActivationLock performs Activation Lock bypass for a device
func (a *BypassActivationLockManagedDeviceAction) bypassActivationLock(ctx context.Context, deviceID string) error {
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		BypassActivationLock().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
