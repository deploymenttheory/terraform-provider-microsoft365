package graphBetaUpdateWindowsDeviceAccount

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type updateResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	email      string
	err        error
}

func (a *UpdateWindowsDeviceAccountAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data UpdateWindowsDeviceAccountActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Updating device accounts for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Windows device account update for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Update device accounts concurrently with error collection
	results := make(chan updateResult, totalDevices)
	var wg sync.WaitGroup

	// Update managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceAccount) {
			defer wg.Done()
			deviceID := d.DeviceID.ValueString()
			email := d.DeviceAccountEmail.ValueString()
			err := a.updateManagedDevice(ctx, d)
			results <- updateResult{deviceID: deviceID, deviceType: "managed", email: email, err: err}
		}(device)
	}

	// Update co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceAccount) {
			defer wg.Done()
			deviceID := d.DeviceID.ValueString()
			email := d.DeviceAccountEmail.ValueString()
			err := a.updateComanagedDevice(ctx, d)
			results <- updateResult{deviceID: deviceID, deviceType: "comanaged", email: email, err: err}
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
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s, %s)", result.deviceID, result.deviceType, result.email))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to update device account on %s device %s (%s): %v",
				result.deviceType, result.deviceID, result.email, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully updated device account on %s device %s (%s)",
				result.deviceType, result.deviceID, result.email))
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
				fmt.Sprintf("Successfully updated device accounts on %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that were successfully updated will apply the new device account configuration. "+
					"The devices may need to be rebooted for changes to take full effect. "+
					"Failed devices may be offline, not Windows collaboration devices, or may have configuration issues.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated device accounts on %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Windows device account update completed on %d device(s). "+
				"Device account credentials, Exchange server settings, and synchronization options have been updated. "+
				"Devices may require a reboot for all changes to take effect. "+
				"Verify functionality after devices restart and reconnect to Exchange and Teams/Skype for Business services.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *UpdateWindowsDeviceAccountAction) updateManagedDevice(ctx context.Context, device ManagedDeviceAccount) error {
	deviceID := device.DeviceID.ValueString()
	email := device.DeviceAccountEmail.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating device account on managed device %s (account: %s)", deviceID, email))

	// Construct the request body
	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		UpdateWindowsDeviceAccount().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to update device account on managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *UpdateWindowsDeviceAccountAction) updateComanagedDevice(ctx context.Context, device ComanagedDeviceAccount) error {
	deviceID := device.DeviceID.ValueString()
	email := device.DeviceAccountEmail.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating device account on co-managed device %s (account: %s)", deviceID, email))

	// Construct the request body
	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		UpdateWindowsDeviceAccount().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to update device account on co-managed device %s: %w", deviceID, err)
	}

	return nil
}
