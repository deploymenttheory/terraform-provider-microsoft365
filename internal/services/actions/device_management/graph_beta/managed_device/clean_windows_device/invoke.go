package graphBetaCleanWindowsManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *CleanWindowsManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data CleanWindowsManagedDeviceActionModel

	tflog.Debug(ctx, "Starting Windows device clean action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for Windows clean", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
		"total_devices":     totalDevices,
	})

	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Windows device clean for %d devices (%d managed, %d co-managed)",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	successCount := 0
	var failedDevices []string
	var lastError error

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		keepUserData := device.KeepUserData.ValueBool()
		err := a.cleanManagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (managed)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to clean managed device", map[string]any{
				"device_id":      deviceID,
				"keep_user_data": keepUserData,
				"error":          err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully initiated clean for managed device", map[string]any{
				"device_id":      deviceID,
				"keep_user_data": keepUserData,
			})
		}

		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		keepUserData := device.KeepUserData.ValueBool()
		err := a.cleanComanagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (comanaged)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to clean co-managed device", map[string]any{
				"device_id":      deviceID,
				"keep_user_data": keepUserData,
				"error":          err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully initiated clean for co-managed device", map[string]any{
				"device_id":      deviceID,
				"keep_user_data": keepUserData,
			})
		}

		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	a.reportResults(ctx, resp, successCount, totalDevices, failedDevices, lastError, ignorePartialFailures)

	tflog.Debug(ctx, "Completed Windows device clean action", map[string]any{
		"action":        ActionName,
		"success_count": successCount,
		"failed_count":  len(failedDevices),
		"total_devices": totalDevices,
	})
}

// cleanManagedDevice performs atomic clean for a managed device
func (a *CleanWindowsManagedDeviceAction) cleanManagedDevice(ctx context.Context, device ManagedDeviceCleanWindows) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Cleaning managed device", map[string]any{"device_id": deviceID})

	requestBody := constructManagedDeviceRequest(ctx, device)
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		CleanWindowsDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated clean for managed device", map[string]any{"device_id": deviceID})
	return nil
}

// cleanComanagedDevice performs atomic clean for a co-managed device
func (a *CleanWindowsManagedDeviceAction) cleanComanagedDevice(ctx context.Context, device ComanagedDeviceCleanWindows) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Cleaning co-managed device", map[string]any{"device_id": deviceID})

	requestBody := constructComanagedDeviceRequest(ctx, device)
	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		CleanWindowsDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated clean for co-managed device", map[string]any{"device_id": deviceID})
	return nil
}

// reportResults handles final result reporting according to ADR-001 principles
func (a *CleanWindowsManagedDeviceAction) reportResults(ctx context.Context, resp *action.InvokeResponse, successCount, totalDevices int, failedDevices []string, lastError error, ignorePartialFailures bool) {
	if len(failedDevices) > 0 {
		if successCount > 0 {
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices processed. %d failed devices ignored. Clean operation initiated for successful devices.",
						successCount, totalDevices, len(failedDevices)),
				})
			} else {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices cleaned. Failed devices: %v",
						successCount, totalDevices, failedDevices),
				})
				resp.Diagnostics.AddWarning(
					"Partial Success",
					fmt.Sprintf("Windows device clean partially completed. %d of %d devices succeeded. Failed devices: %v",
						successCount, totalDevices, failedDevices),
				)
			}
		} else {
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("All %d devices failed clean. Failures ignored per configuration.",
						totalDevices),
				})
			} else {
				errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
				return
			}
		}
	} else {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Windows device clean initiated successfully for all %d devices. "+
				"Devices will begin the clean process. This may take several minutes to complete.",
				successCount),
		})
	}
}
