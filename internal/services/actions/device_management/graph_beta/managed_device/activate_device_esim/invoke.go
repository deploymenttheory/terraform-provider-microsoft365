package graphBetaActivateDeviceEsimManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ActivateDeviceEsimManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ActivateDeviceEsimManagedDeviceActionModel

	tflog.Debug(ctx, "Starting eSIM activation action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for eSIM activation", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
		"total_devices":     totalDevices,
	})

	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting eSIM activation for %d devices (%d managed, %d co-managed)",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	successCount := 0
	var failedDevices []string
	var lastError error

	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.activateEsimManagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (managed)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to activate eSIM for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully activated eSIM for managed device", map[string]any{
				"device_id": deviceID,
			})
		}

		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.activateEsimComanagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (comanaged)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to activate eSIM for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully activated eSIM for co-managed device", map[string]any{
				"device_id": deviceID,
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

	tflog.Debug(ctx, "Completed eSIM activation action", map[string]any{
		"action":        ActionName,
		"success_count": successCount,
		"failed_count":  len(failedDevices),
		"total_devices": totalDevices,
	})
}

// activateEsimManagedDevice performs atomic eSIM activation for a managed device
func (a *ActivateDeviceEsimManagedDeviceAction) activateEsimManagedDevice(ctx context.Context, device ManagedDeviceActivateEsim) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Activating eSIM for managed device", map[string]any{"device_id": deviceID})

	requestBody := constructManagedDeviceRequest(ctx, device)
	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		ActivateDeviceEsim().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated eSIM activation for managed device", map[string]any{"device_id": deviceID})
	return nil
}

// activateEsimComanagedDevice performs atomic eSIM activation for a co-managed device
func (a *ActivateDeviceEsimManagedDeviceAction) activateEsimComanagedDevice(ctx context.Context, device ComanagedDeviceActivateEsim) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Activating eSIM for co-managed device", map[string]any{"device_id": deviceID})

	requestBody := constructComanagedDeviceRequest(ctx, device)
	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		ActivateDeviceEsim().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated eSIM activation for co-managed device", map[string]any{"device_id": deviceID})
	return nil
}

// reportResults handles final result reporting according to ADR-001 principles
func (a *ActivateDeviceEsimManagedDeviceAction) reportResults(ctx context.Context, resp *action.InvokeResponse, successCount, totalDevices int, failedDevices []string, lastError error, ignorePartialFailures bool) {
	if len(failedDevices) > 0 {
		if successCount > 0 {
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices processed. %d failed devices ignored. Successfully activated devices will receive eSIM configuration.",
						successCount, totalDevices, len(failedDevices)),
				})
			} else {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices activated. Failed devices: %v",
						successCount, totalDevices, failedDevices),
				})
				resp.Diagnostics.AddWarning(
					"Partial Success",
					fmt.Sprintf("eSIM activation partially completed. %d of %d devices succeeded. Failed devices: %v",
						successCount, totalDevices, failedDevices),
				)
			}
		} else {
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("All %d devices failed activation. Failures ignored per configuration.",
						totalDevices),
				})
			} else {
				errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
				return
			}
		}
	} else {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("eSIM activation completed successfully for all %d devices.",
				successCount),
		})
	}
}
