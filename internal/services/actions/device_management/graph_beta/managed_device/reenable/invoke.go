package graphBetaReenableManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type reenableResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *ReenableManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ReenableManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert framework lists to Go slices
	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)
	tflog.Debug(ctx, fmt.Sprintf("Re-enabling %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting re-enable for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Re-enable devices concurrently with error collection
	results := make(chan reenableResult, totalDevices)
	var wg sync.WaitGroup

	// Re-enable managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.reenableManagedDevice(ctx, id)
			results <- reenableResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Re-enable co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.reenableComanagedDevice(ctx, id)
			results <- reenableResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to re-enable %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully re-enabled %s device %s",
				result.deviceType, result.deviceID))

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Re-enabled", result.deviceID, result.deviceType),
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
				fmt.Sprintf("Successfully re-enabled %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Re-enabled devices will now be able to sync with Intune and receive policy updates.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully re-enabled %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Re-enable complete: %d device(s) have been re-enabled for Intune management.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *ReenableManagedDeviceAction) reenableManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Re-enabling managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		Reenable().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to re-enable managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *ReenableManagedDeviceAction) reenableComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Re-enabling co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		Reenable().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to re-enable co-managed device %s: %w", deviceID, err)
	}

	return nil
}
