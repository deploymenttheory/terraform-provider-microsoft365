package graphBetaResetManagedDevicePasscode

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resetResult struct {
	deviceID string
	err      error
}

func (a *ResetManagedDevicePasscodeAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ResetManagedDevicePasscodeActionModel

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
		Message: fmt.Sprintf("Starting passcode reset for %d managed device(s)...", totalDevices),
	})

	// Reset passcodes concurrently with error collection
	results := make(chan resetResult, totalDevices)
	var wg sync.WaitGroup

	for _, deviceID := range deviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.resetPasscode(ctx, id)
			results <- resetResult{deviceID: id, err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to reset passcode for device %s: %v", result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully reset passcode for device %s", result.deviceID))
		}

		// Send progress update
		progress := float64(successCount+len(failedDevices)) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Reset %d of %d device passcodes (%.0f%% complete)", successCount+len(failedDevices), totalDevices, progress),
		})
	}

	// Report results
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully reset passcodes for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"IMPORTANT: Check the Intune portal for the new temporary passcodes for the successfully reset devices. "+
					"Communicate these passcodes to the device users securely.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully reset passcodes for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Passcode reset complete: %d device(s) successfully reset. "+
				"IMPORTANT: Check the Intune portal (Devices > All devices > select device > Reset passcode) "+
				"to retrieve the new temporary passcodes and communicate them to device users.", successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *ResetManagedDevicePasscodeAction) resetPasscode(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Resetting passcode for device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		ResetPasscode().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to reset passcode for device %s: %w", deviceID, err)
	}

	return nil
}
