package graphBetaEnableLostModeManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type enableLostModeResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *EnableLostModeManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data EnableLostModeManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Enabling lost mode for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting enable lost mode for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Enable lost mode on devices concurrently with error collection
	results := make(chan enableLostModeResult, totalDevices)
	var wg sync.WaitGroup

	// Enable lost mode on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceLostMode) {
			defer wg.Done()
			err := a.enableLostModeManagedDevice(ctx, d)
			results <- enableLostModeResult{deviceID: d.DeviceID.ValueString(), deviceType: "managed", err: err}
		}(device)
	}

	// Enable lost mode on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceLostMode) {
			defer wg.Done()
			err := a.enableLostModeComanagedDevice(ctx, d)
			results <- enableLostModeResult{deviceID: d.DeviceID.ValueString(), deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to enable lost mode for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully enabled lost mode for %s device %s",
				result.deviceType, result.deviceID))
		}

		// Send progress update
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
				fmt.Sprintf("Successfully enabled lost mode for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that had lost mode enabled are now locked with the custom message displayed.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully enabled lost mode for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Enable lost mode complete: %d device(s) successfully locked and secured. "+
				"Devices are now in lost mode with the lock screen message displayed and location tracking enabled.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *EnableLostModeManagedDeviceAction) enableLostModeManagedDevice(ctx context.Context, device ManagedDeviceLostMode) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Enabling lost mode for managed device with ID: %s", deviceID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		EnableLostMode().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to enable lost mode for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *EnableLostModeManagedDeviceAction) enableLostModeComanagedDevice(ctx context.Context, device ComanagedDeviceLostMode) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Enabling lost mode for co-managed device with ID: %s", deviceID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		EnableLostMode().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to enable lost mode for co-managed device %s: %w", deviceID, err)
	}

	return nil
}
