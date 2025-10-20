package graphBetaLocateManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type locateDeviceResult struct {
	deviceID string
	err      error
}

func (a *LocateManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data LocateManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

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
	tflog.Debug(ctx, fmt.Sprintf("Performing action %s for %d device(s)", ActionName, totalDevices))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting locate device request for %d managed device(s)...", totalDevices),
	})

	// Locate devices concurrently with error collection
	results := make(chan locateDeviceResult, totalDevices)
	var wg sync.WaitGroup

	for _, deviceID := range deviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.locateDevice(ctx, id)
			results <- locateDeviceResult{deviceID: id, err: err}
		}(deviceID)
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
			failedDevices = append(failedDevices, result.deviceID)
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to send locate request to device %s: %v", result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully sent locate request to device %s", result.deviceID))
		}

		// Send progress update
		progress := float64(successCount+len(failedDevices)) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)", successCount+len(failedDevices), totalDevices, progress),
		})
	}

	// Report results
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully sent locate request to %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that received the locate command will report their location when online. "+
					"Location data will be available in the Microsoft Intune admin center.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully sent locate request to %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Locate device complete: %d device(s) received location request. "+
				"Devices will report their location when online and location services are enabled. "+
				"View location data in the Microsoft Intune admin center under device properties.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *LocateManagedDeviceAction) locateDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Sending locate device request to device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		LocateDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to send locate request to device %s: %w", deviceID, err)
	}

	return nil
}
