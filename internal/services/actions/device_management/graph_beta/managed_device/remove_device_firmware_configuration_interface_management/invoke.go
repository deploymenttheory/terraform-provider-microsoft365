package graphBetaRemoveDeviceFirmwareConfigurationInterfaceManagementManagedDevice

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

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceActionModel

	tflog.Debug(ctx, "Starting remove DFCI management action", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for DFCI removal", map[string]any{
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
			Warning(validationResult.UnsupportedManagedDevices, "managed device", "are not Windows devices. DFCI is only supported on Windows devices with compatible firmware").
			Warning(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Windows devices. DFCI is only supported on Windows devices with compatible firmware")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("DFCI management removal", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)))

	// Process managed devices sequentially
	for _, deviceID := range managedDeviceIDs {
		err := a.removeDFCIManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to remove DFCI management from managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("DFCI management removed")
			tflog.Info(ctx, "Successfully removed DFCI management from managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, deviceID := range comanagedDeviceIDs {
		err := a.removeDFCIComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(err.Error())
			tflog.Error(ctx, "Failed to remove DFCI management from co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("DFCI management removed")
			tflog.Info(ctx, "Successfully removed DFCI management from co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("DFCI management removal")
			tflog.Warn(ctx, "DFCI removal action completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("DFCI Management Removal Failed", "remove DFCI management from devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("removed DFCI management from all devices. These devices will no longer support remote UEFI/BIOS configuration via Intune. Standard MDM management continues")
	}

	tflog.Info(ctx, "DFCI removal action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) removeDFCIManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Removing DFCI management from managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RemoveDeviceFirmwareConfigurationInterfaceManagement().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) removeDFCIComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Removing DFCI management from co-managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RemoveDeviceFirmwareConfigurationInterfaceManagement().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
