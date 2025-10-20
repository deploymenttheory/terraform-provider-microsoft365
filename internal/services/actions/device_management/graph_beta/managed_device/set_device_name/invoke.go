package graphBetaSetDeviceNameManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type setDeviceNameResult struct {
	deviceID   string
	deviceName string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *SetDeviceNameManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SetDeviceNameManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Setting device names for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting set device name for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Set device names concurrently with error collection
	results := make(chan setDeviceNameResult, totalDevices)
	var wg sync.WaitGroup

	// Set names on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceSetName) {
			defer wg.Done()
			err := a.setDeviceNameManagedDevice(ctx, d)
			results <- setDeviceNameResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceName: d.DeviceName.ValueString(),
				deviceType: "managed",
				err:        err,
			}
		}(device)
	}

	// Set names on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceSetName) {
			defer wg.Done()
			err := a.setDeviceNameComanagedDevice(ctx, d)
			results <- setDeviceNameResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceName: d.DeviceName.ValueString(),
				deviceType: "comanaged",
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
			tflog.Error(ctx, fmt.Sprintf("Failed to set device name for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully set device name to '%s' for %s device %s",
				result.deviceName, result.deviceType, result.deviceID))
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
				fmt.Sprintf("Successfully set device names for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that received the command will have their names updated after their next check-in.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully set device names for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Set device name complete: %d device(s) will have their names updated after next check-in.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *SetDeviceNameManagedDeviceAction) setDeviceNameManagedDevice(ctx context.Context, device ManagedDeviceSetName) error {
	deviceID := device.DeviceID.ValueString()
	deviceName := device.DeviceName.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Setting device name to '%s' for managed device with ID: %s", deviceName, deviceID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SetDeviceName().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to set device name for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *SetDeviceNameManagedDeviceAction) setDeviceNameComanagedDevice(ctx context.Context, device ComanagedDeviceSetName) error {
	deviceID := device.DeviceID.ValueString()
	deviceName := device.DeviceName.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Setting device name to '%s' for co-managed device with ID: %s", deviceName, deviceID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SetDeviceName().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to set device name for co-managed device %s: %w", deviceID, err)
	}

	return nil
}

