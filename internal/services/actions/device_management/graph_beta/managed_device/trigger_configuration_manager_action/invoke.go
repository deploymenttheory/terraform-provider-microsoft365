package graphBetaTriggerConfigurationManagerActionManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type configManagerActionResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	action     string
	err        error
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TriggerConfigurationManagerActionManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Triggering Configuration Manager actions for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Configuration Manager action triggers for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Trigger actions concurrently with error collection
	results := make(chan configManagerActionResult, totalDevices)
	var wg sync.WaitGroup

	// Trigger actions on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceConfigManagerAction) {
			defer wg.Done()
			err := a.triggerConfigManagerActionManagedDevice(ctx, d)
			results <- configManagerActionResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "managed",
				action:     d.Action.ValueString(),
				err:        err,
			}
		}(device)
	}

	// Trigger actions on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceConfigManagerAction) {
			defer wg.Done()
			err := a.triggerConfigManagerActionComanagedDevice(ctx, d)
			results <- configManagerActionResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "comanaged",
				action:     d.Action.ValueString(),
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
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s, action: %s)", result.deviceID, result.deviceType, result.action))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to trigger Configuration Manager action '%s' for %s device %s: %v",
				result.action, result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully triggered Configuration Manager action '%s' for %s device %s",
				result.action, result.deviceType, result.deviceID))

			// Display success message
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Triggered action '%s'",
					result.deviceID, result.deviceType, result.action),
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
				fmt.Sprintf("Successfully triggered Configuration Manager actions for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that received the command will execute the requested Configuration Manager action.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully triggered Configuration Manager actions for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Configuration Manager action triggers complete: %d device(s) will execute the requested actions.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) triggerConfigManagerActionManagedDevice(ctx context.Context, device ManagedDeviceConfigManagerAction) error {
	deviceID := device.DeviceID.ValueString()
	action := device.Action.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Triggering Configuration Manager action '%s' for managed device with ID: %s", action, deviceID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		TriggerConfigurationManagerAction().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to trigger Configuration Manager action '%s' for managed device %s: %w", action, deviceID, err)
	}

	return nil
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) triggerConfigManagerActionComanagedDevice(ctx context.Context, device ComanagedDeviceConfigManagerAction) error {
	deviceID := device.DeviceID.ValueString()
	action := device.Action.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Triggering Configuration Manager action '%s' for co-managed device with ID: %s", action, deviceID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		TriggerConfigurationManagerAction().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to trigger Configuration Manager action '%s' for co-managed device %s: %w", action, deviceID, err)
	}

	return nil
}
