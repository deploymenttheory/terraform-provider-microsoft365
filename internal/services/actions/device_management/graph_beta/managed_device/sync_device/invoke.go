package graphBetaSyncManagedDevice

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

func (a *SyncManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SyncManagedDeviceActionModel

	tflog.Debug(ctx, "Starting device sync action", map[string]any{"action": ActionName})

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
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)
	tflog.Debug(ctx, "Processing devices for sync", map[string]any{
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not supported for device sync. Sync is supported on Windows, macOS, iOS, iPadOS, and Android devices only").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not supported for device sync. Sync is supported on Windows, macOS, iOS, iPadOS, and Android devices only")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("device sync", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)))

	// Process managed devices sequentially
	for _, deviceID := range managedDeviceIDs {
		err := a.syncManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("sync command failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to sync managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("sync command sent")
			tflog.Info(ctx, "Successfully sent sync command to managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices sequentially
	for _, deviceID := range comanagedDeviceIDs {
		err := a.syncComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("sync command failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to sync co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("sync command sent")
			tflog.Info(ctx, "Successfully sent sync command to co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("device sync")
			tflog.Warn(ctx, "Device sync action completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Device Sync Failed", "send sync commands to devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("sent sync commands to all devices. Online devices will check in with Intune immediately (within 1-5 minutes) and apply the latest policies, configurations, and updates. Offline devices will sync when they come back online. Check sync status in the Microsoft Intune admin center")
	}

	tflog.Info(ctx, "Device sync action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *SyncManagedDeviceAction) syncManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Sending sync command to managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SyncDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *SyncManagedDeviceAction) syncComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Sending sync command to co-managed device", map[string]any{
		"device_id": deviceID,
	})

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SyncDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
