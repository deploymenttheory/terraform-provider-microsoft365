package graphBetaPlayLostModeSoundManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type playLostModeSoundResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *PlayLostModeSoundManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data PlayLostModeSoundManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Playing lost mode sound for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting play lost mode sound for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Play sound on devices concurrently with error collection
	results := make(chan playLostModeSoundResult, totalDevices)
	var wg sync.WaitGroup

	// Play sound on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDevicePlaySound) {
			defer wg.Done()
			err := a.playLostModeSoundManagedDevice(ctx, d)
			results <- playLostModeSoundResult{deviceID: d.DeviceID.ValueString(), deviceType: "managed", err: err}
		}(device)
	}

	// Play sound on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDevicePlaySound) {
			defer wg.Done()
			err := a.playLostModeSoundComanagedDevice(ctx, d)
			results <- playLostModeSoundResult{deviceID: d.DeviceID.ValueString(), deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to play lost mode sound for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully played lost mode sound for %s device %s",
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
				fmt.Sprintf("Successfully played lost mode sound for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that received the command will play the lost mode sound to assist in locating them.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully played lost mode sound for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Play lost mode sound complete: %d device(s) will now play an audible alert. "+
				"The sound will play even if the device is in silent mode, helping to locate the device physically.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *PlayLostModeSoundManagedDeviceAction) playLostModeSoundManagedDevice(ctx context.Context, device ManagedDevicePlaySound) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Playing lost mode sound for managed device with ID: %s", deviceID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		PlayLostModeSound().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to play lost mode sound for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *PlayLostModeSoundManagedDeviceAction) playLostModeSoundComanagedDevice(ctx context.Context, device ComanagedDevicePlaySound) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Playing lost mode sound for co-managed device with ID: %s", deviceID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		PlayLostModeSound().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to play lost mode sound for co-managed device %s: %w", deviceID, err)
	}

	return nil
}
