package graphBetaSyncManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type syncResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *SyncManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SyncManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)
	tflog.Debug(ctx, fmt.Sprintf("Performing action %s for %d managed device(s) and %d co-managed device(s)",
		ActionName, len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting device sync for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Sync devices concurrently with error collection
	results := make(chan syncResult, totalDevices)
	var wg sync.WaitGroup

	// Sync managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.syncManagedDevice(ctx, id)
			results <- syncResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Sync co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.syncComanagedDevice(ctx, id)
			results <- syncResult{deviceID: id, deviceType: "comanaged", err: err}
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
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s)", result.deviceID, result.deviceType))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to sync %s device %s: %v", result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully synced %s device %s", result.deviceType, result.deviceID))
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
				fmt.Sprintf("Successfully synced %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that were successfully synced will check in with Intune and apply the latest policies and configurations. "+
					"Failed devices may be offline, have network issues, or may not exist.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully synced %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Sync complete: %d device(s) will check in with Intune and apply the latest policies, "+
				"configurations, and updates. Check-in typically completes within 1-5 minutes for online devices. "+
				"View sync status in the Microsoft Intune admin center.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *SyncManagedDeviceAction) syncManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Syncing managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SyncDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to sync managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *SyncManagedDeviceAction) syncComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Syncing co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SyncDevice().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to sync co-managed device %s: %w", deviceID, err)
	}

	return nil
}
