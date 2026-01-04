package graphBetaLocateManagedDevice

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

func (a *LocateManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data LocateManagedDeviceActionModel

	tflog.Debug(ctx, "Starting locate device action", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for location", map[string]any{
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
			Warning(validationResult.OfflineDevices, "device", "may be offline or not properly registered. The locate request will be queued and processed when the device comes online").
			Warning(validationResult.LocationDisabledDevices, "device", "may have limited location capabilities. For iOS/iPadOS devices, supervision provides more reliable location reporting").
			Warning(validationResult.UnsupportedOSDevices, "device", "have operating systems that do not support locate device (must be Windows, iOS, iPadOS, or Android)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("device location", fmt.Sprintf("%d devices", totalDevices))

	// Process devices sequentially
	for _, deviceID := range deviceIDs {
		err := a.locateDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to locate device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("location request sent successfully")
			tflog.Info(ctx, "Successfully sent location request to device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("device location")
			tflog.Warn(ctx, "Device location completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Device Location Failed", "send location requests to devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("sent location requests to all devices. Location data will appear in Intune once devices report their position")
	}

	tflog.Info(ctx, "Locate device action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *LocateManagedDeviceAction) locateDevice(ctx context.Context, deviceID string) error {
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		LocateDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
