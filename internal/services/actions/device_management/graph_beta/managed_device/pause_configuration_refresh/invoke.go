package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type pauseConfigResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	pauseTime  int64
	err        error
}

func (a *PauseConfigurationRefreshManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data PauseConfigurationRefreshManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Pausing configuration refresh for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting configuration refresh pause for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Pause configuration refresh concurrently with error collection
	results := make(chan pauseConfigResult, totalDevices)
	var wg sync.WaitGroup

	// Pause on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDevicePauseConfig) {
			defer wg.Done()
			err := a.pauseConfigManagedDevice(ctx, d)
			results <- pauseConfigResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "managed",
				pauseTime:  d.PauseTimePeriodInMinutes.ValueInt64(),
				err:        err,
			}
		}(device)
	}

	// Pause on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDevicePauseConfig) {
			defer wg.Done()
			err := a.pauseConfigComanagedDevice(ctx, d)
			results <- pauseConfigResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "comanaged",
				pauseTime:  d.PauseTimePeriodInMinutes.ValueInt64(),
				err:        err,
			}
		}(device)
	}

	// Close results channel once all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and track progress
	successCount := 0
	var failedDevices []string
	var lastError error

	for result := range results {
		if result.err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s)", result.deviceID, result.deviceType))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to pause configuration refresh for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully paused configuration refresh for %s device %s for %d minutes",
				result.deviceType, result.deviceID, result.pauseTime))

			// Calculate hours and minutes for display
			hours := result.pauseTime / 60
			minutes := result.pauseTime % 60
			timeDesc := ""
			if hours > 0 && minutes > 0 {
				timeDesc = fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
			} else if hours > 0 {
				timeDesc = fmt.Sprintf("%d hour(s)", hours)
			} else {
				timeDesc = fmt.Sprintf("%d minute(s)", minutes)
			}

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Configuration refresh paused for %s", result.deviceID, result.deviceType, timeDesc),
			})
		}

		// Send overall progress update
		progress := float64(successCount+len(failedDevices)) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				successCount+len(failedDevices), totalDevices, progress),
		})
	}

	// Report results
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully paused configuration refresh for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices with successful pauses will not receive new policy updates during their pause periods. "+
					"Configuration refresh will automatically resume after the specified time.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully paused configuration refresh for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ Configuration refresh pause complete: %d device(s) will not receive policy updates during their pause periods.\n\n"+
				"Important reminders:\n"+
				"- Configuration refresh will automatically resume after the pause period\n"+
				"- Existing policies remain in effect\n"+
				"- Users can still manually sync from Company Portal\n"+
				"- Critical security updates may still be applied\n"+
				"- Monitor devices during the pause period",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *PauseConfigurationRefreshManagedDeviceAction) pauseConfigManagedDevice(ctx context.Context, device ManagedDevicePauseConfig) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Pausing configuration refresh for managed device with ID: %s for %d minutes", deviceID, device.PauseTimePeriodInMinutes.ValueInt64()))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		PauseConfigurationRefresh().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to pause configuration refresh for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *PauseConfigurationRefreshManagedDeviceAction) pauseConfigComanagedDevice(ctx context.Context, device ComanagedDevicePauseConfig) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Pausing configuration refresh for co-managed device with ID: %s for %d minutes", deviceID, device.PauseTimePeriodInMinutes.ValueInt64()))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		PauseConfigurationRefresh().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to pause configuration refresh for co-managed device %s: %w", deviceID, err)
	}

	return nil
}
