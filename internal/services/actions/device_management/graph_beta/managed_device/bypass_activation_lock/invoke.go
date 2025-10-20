package graphBetaBypassActivationLockManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *BypassActivationLockManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data BypassActivationLockManagedDeviceActionModel

	tflog.Debug(ctx, "Starting Activation Lock bypass action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, "Processing devices for Activation Lock bypass", map[string]any{
		"total_devices": totalDevices,
	})

	// Get configuration values with defaults
	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Activation Lock bypass for %d devices", totalDevices),
	})

	// Process devices atomically - no retry logic, single attempt per device
	successCount := 0
	var failedDevices []string
	var lastError error

	// Process each device sequentially
	for _, deviceID := range deviceIDs {
		err := a.bypassActivationLock(ctx, deviceID)
		if err != nil {
			failedDevices = append(failedDevices, deviceID)
			lastError = err
			tflog.Error(ctx, "Failed to bypass Activation Lock for device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully bypassed Activation Lock for device", map[string]any{
				"device_id": deviceID,
			})
		}

		// Send progress update
		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	// Report final results using ADR-001 principles
	a.reportResults(ctx, resp, successCount, totalDevices, failedDevices, lastError, ignorePartialFailures)

	tflog.Debug(ctx, "Completed Activation Lock bypass action", map[string]any{
		"action":        ActionName,
		"success_count": successCount,
		"failed_count":  len(failedDevices),
		"total_devices": totalDevices,
	})
}

// bypassActivationLock performs atomic Activation Lock bypass for a device
func (a *BypassActivationLockManagedDeviceAction) bypassActivationLock(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, "Bypassing Activation Lock for device", map[string]any{"device_id": deviceID})

	if err := constructBypassActivationLockRequest(ctx); err != nil {
		return err
	}

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		BypassActivationLock().
		Post(ctx, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated Activation Lock bypass for device", map[string]any{"device_id": deviceID})
	return nil
}

// reportResults handles final result reporting according to ADR-001 principles
func (a *BypassActivationLockManagedDeviceAction) reportResults(ctx context.Context, resp *action.InvokeResponse, successCount, totalDevices int, failedDevices []string, lastError error, ignorePartialFailures bool) {
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices processed. %d failed devices ignored. Bypass codes generated for successful devices.",
						successCount, totalDevices, len(failedDevices)),
				})
			} else {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("Partial success: %d of %d devices processed. Failed devices: %v",
						successCount, totalDevices, failedDevices),
				})
				// For partial failures without ignore flag, add diagnostic warning
				resp.Diagnostics.AddWarning(
					"Partial Success",
					fmt.Sprintf("Activation Lock bypass partially completed. %d of %d devices succeeded. Failed devices: %v. "+
						"Bypass codes for successful devices are available in Intune device properties.",
						successCount, totalDevices, failedDevices),
				)
			}
		} else {
			// Complete failure
			if ignorePartialFailures {
				resp.SendProgress(action.InvokeProgressEvent{
					Message: fmt.Sprintf("All %d devices failed bypass. Failures ignored per configuration.",
						totalDevices),
				})
			} else {
				// Complete failure should be handled as an error
				errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
				return
			}
		}
	} else {
		// Complete success
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Activation Lock bypass completed successfully for all %d devices. "+
				"Bypass codes are now available in Intune device properties.",
				successCount),
		})
	}
}
