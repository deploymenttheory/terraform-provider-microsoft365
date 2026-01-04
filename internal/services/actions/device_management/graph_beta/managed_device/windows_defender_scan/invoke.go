package graphBetaWindowsDefenderScan

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

func (a *WindowsDefenderScanAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data WindowsDefenderScanActionModel

	tflog.Debug(ctx, "Starting Windows Defender scan", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for Windows Defender scan", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
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

		validationResult, err := validateRequest(ctx, a.client, data.ManagedDevices, data.ComanagedDevices)
		if err != nil {
			tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError("Device Validation Failed", fmt.Sprintf("Failed to validate devices: %s", err.Error()))
			return
		}

		results := validation.NewResults().
			Error(validationResult.NonExistentManagedDevices, "managed device", "do not exist or are not managed by Intune").
			Error(validationResult.NonWindowsManagedDevices, "managed device", "are not Windows devices. Windows Defender scan only works on Windows devices").
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Error(validationResult.NonWindowsComanagedDevices, "co-managed device", "are not Windows devices. Windows Defender scan only works on Windows devices")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("Windows Defender scans", fmt.Sprintf("%d devices (%d managed, %d co-managed)", totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)))

	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		scanType := "full"
		if device.QuickScan.ValueBool() {
			scanType = "quick"
		}

		err := a.scanManagedDevice(ctx, deviceID, device.QuickScan.ValueBool())
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("%s scan failed: %s", scanType, err.Error()))
			tflog.Error(ctx, "Failed to initiate scan on managed device", map[string]any{"device_id": deviceID, "scan_type": scanType, "error": err.Error()})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("%s scan initiated", scanType))
			tflog.Info(ctx, "Successfully initiated scan on managed device", map[string]any{"device_id": deviceID, "scan_type": scanType})
		}
	}

	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		scanType := "full"
		if device.QuickScan.ValueBool() {
			scanType = "quick"
		}

		err := a.scanComanagedDevice(ctx, deviceID, device.QuickScan.ValueBool())
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("%s scan failed: %s", scanType, err.Error()))
			tflog.Error(ctx, "Failed to initiate scan on co-managed device", map[string]any{"device_id": deviceID, "scan_type": scanType, "error": err.Error()})
		} else {
			progressTracker.Device(deviceID, "").Succeeded(fmt.Sprintf("%s scan initiated", scanType))
			tflog.Info(ctx, "Successfully initiated scan on co-managed device", map[string]any{"device_id": deviceID, "scan_type": scanType})
		}
	}

	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("Windows Defender scans")
			tflog.Warn(ctx, "Windows Defender scan completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Windows Defender Scan Failed", "initiate Windows Defender scans on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("initiated Windows Defender scans on all devices. Scans will begin immediately on online devices. Quick scans take 5-15 minutes, full scans take 30+ minutes to hours. Results will be reported to Microsoft Intune admin center. Threats found will be quarantined automatically")
	}

	tflog.Info(ctx, "Windows Defender scan completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *WindowsDefenderScanAction) scanManagedDevice(ctx context.Context, deviceID string, quickScan bool) error {
	scanType := "full"
	if quickScan {
		scanType = "quick"
	}
	tflog.Debug(ctx, "Initiating Windows Defender scan on managed device", map[string]any{"device_id": deviceID, "scan_type": scanType})

	requestBody := constructManagedDeviceRequest(ctx, quickScan)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderScan().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *WindowsDefenderScanAction) scanComanagedDevice(ctx context.Context, deviceID string, quickScan bool) error {
	scanType := "full"
	if quickScan {
		scanType = "quick"
	}
	tflog.Debug(ctx, "Initiating Windows Defender scan on co-managed device", map[string]any{"device_id": deviceID, "scan_type": scanType})

	requestBody := constructComanagedDeviceRequest(ctx, quickScan)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderScan().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
