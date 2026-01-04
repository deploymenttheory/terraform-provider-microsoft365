package graphBetaLogoutSharedAppleDeviceActiveUser

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

func (a *LogoutSharedAppleDeviceActiveUserAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data LogoutSharedAppleDeviceActiveUserActionModel

	tflog.Debug(ctx, "Starting logout shared Apple device active user action", map[string]any{"action": ActionName})

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

	// Convert framework list to Go slice
	var deviceIDs []string
	if !data.DeviceIDs.IsNull() && !data.DeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, "Processing Shared iPad devices for active user logout", map[string]any{
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
			Error(validationResult.NonIPadDevices, "device", "are not iPadOS devices (Shared iPad mode requires iPadOS)").
			Warning(validationResult.UnsupervisedDevices, "device", "are not supervised (Shared iPad mode requires supervised devices enrolled via DEP/ABM)").
			Warning(validationResult.PotentiallyNotSharedDevices, "device", "may not support Shared iPad mode (likely iPhones, not iPads)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("active user logout", fmt.Sprintf("%d Shared iPad devices", totalDevices))

	// Process devices sequentially
	for _, deviceID := range deviceIDs {
		err := a.logoutActiveUser(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to logout active user from Shared iPad", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("active user logged out successfully")
			tflog.Info(ctx, "Successfully logged out active user from Shared iPad", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("active user logout")
			tflog.Warn(ctx, "Active user logout completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Active User Logout Failed", "logout active users from Shared iPad devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("logged out active users from all Shared iPad devices. Devices have returned to the login screen")
	}

	tflog.Info(ctx, "Logout shared Apple device active user action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *LogoutSharedAppleDeviceActiveUserAction) logoutActiveUser(ctx context.Context, deviceID string) error {
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		LogoutSharedAppleDeviceActiveUser().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
