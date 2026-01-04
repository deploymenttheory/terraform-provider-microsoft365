package graphBetaRecoverManagedDevicePasscode

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

func (a *RecoverManagedDevicePasscodeAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RecoverManagedDevicePasscodeActionModel

	tflog.Debug(ctx, "Starting recover passcode action", map[string]any{"action": ActionName})

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

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, "Processing devices for passcode recovery", map[string]any{
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

		validationResult, err := validateRequest(ctx, a.client, deviceIDs)
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
			Warning(validationResult.UnsupportedOSDevices, "device", "may have limited or no support for passcode recovery (primarily supported on supervised iOS/iPadOS devices)").
			Warning(validationResult.UnsupervisedDevices, "iOS/iPadOS device", "are not supervised (passcode recovery works best with supervised devices enrolled via DEP/ABM)")

		if results.Report(resp) {
			return
		}

		// Add general escrow warning if we have devices
		if len(deviceIDs) > 0 {
			resp.Diagnostics.AddWarning(
				"Passcode Escrow Requirement",
				fmt.Sprintf("Passcode recovery requires that passcodes were escrowed during device enrollment. "+
					"If passcodes were not escrowed for the %d device(s) in this action, recovery will fail. "+
					"Check device enrollment profiles to ensure passcode escrow is enabled. "+
					"If recovery fails, consider using the reset passcode action instead, which generates a new temporary passcode.",
					len(deviceIDs)),
			)
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("passcode recovery", fmt.Sprintf("%d devices", totalDevices))

	// Process devices sequentially
	for _, deviceID := range deviceIDs {
		err := a.recoverPasscode(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to recover passcode for device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("passcode recovered")
			tflog.Info(ctx, "Successfully recovered passcode for device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("passcode recovery")
			tflog.Warn(ctx, "Passcode recovery completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Passcode Recovery Failed", "recover passcodes for devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("recovered passcodes for all devices. Retrieved passcodes are available in the Microsoft Intune admin center under device properties. Securely communicate passcodes to authorized users")
	}

	tflog.Info(ctx, "Recover passcode action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *RecoverManagedDevicePasscodeAction) recoverPasscode(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Recovering passcode for device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RecoverPasscode().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
