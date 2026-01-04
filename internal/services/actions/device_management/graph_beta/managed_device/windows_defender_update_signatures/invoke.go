package graphBetaWindowsDefenderUpdateSignatures

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

func (a *WindowsDefenderUpdateSignaturesAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data WindowsDefenderUpdateSignaturesActionModel

	tflog.Debug(ctx, "Starting Windows Defender signature update", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	tflog.Debug(ctx, "Processing devices for signature update", map[string]any{
		"managed_devices":   len(managedDeviceIDs),
		"comanaged_devices": len(comanagedDeviceIDs),
		"total_devices":     totalDevices,
	})

	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	validateDeviceExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateDeviceExists = data.ValidateDeviceExists.ValueBool()
	}

	if validateDeviceExists {
		tflog.Debug(ctx, "Performing device validation via API")

		validationResult, err := validateRequest(ctx, a.client, managedDeviceIDs, comanagedDeviceIDs)
		if err != nil {
			tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError("Device Validation Failed", fmt.Sprintf("Failed to validate devices: %s", err.Error()))
			return
		}

		results := validation.NewResults().
			Error(validationResult.NonExistentManagedDevices, "managed device", "do not exist or are not managed by Intune").
			Error(validationResult.NonWindowsManagedDevices, "managed device", "are not Windows devices. Windows Defender signature updates only work on Windows devices").
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Error(validationResult.NonWindowsComanagedDevices, "co-managed device", "are not Windows devices. Windows Defender signature updates only work on Windows devices")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("Windows Defender signature updates", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)))

	for _, deviceID := range managedDeviceIDs {
		err := a.updateManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("signature update failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to update signatures on managed device", map[string]any{"device_id": deviceID, "error": err.Error()})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("signature update initiated")
			tflog.Info(ctx, "Successfully initiated signature update on managed device", map[string]any{"device_id": deviceID})
		}
	}

	for _, deviceID := range comanagedDeviceIDs {
		err := a.updateComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("signature update failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to update signatures on co-managed device", map[string]any{"device_id": deviceID, "error": err.Error()})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("signature update initiated")
			tflog.Info(ctx, "Successfully initiated signature update on co-managed device", map[string]any{"device_id": deviceID})
		}
	}

	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("Windows Defender signature updates")
			tflog.Warn(ctx, "Signature update completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Windows Defender Signature Update Failed", "update Windows Defender signatures on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("initiated Windows Defender signature updates on all devices. Devices will download the latest threat definitions and update their antivirus protection. Updates complete within 1-5 minutes for online devices with internet connectivity. View update status in the Microsoft Intune admin center")
	}

	tflog.Info(ctx, "Signature update completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *WindowsDefenderUpdateSignaturesAction) updateManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Updating Windows Defender signatures on managed device", map[string]any{"device_id": deviceID})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderUpdateSignatures().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *WindowsDefenderUpdateSignaturesAction) updateComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Updating Windows Defender signatures on co-managed device", map[string]any{"device_id": deviceID})

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderUpdateSignatures().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
