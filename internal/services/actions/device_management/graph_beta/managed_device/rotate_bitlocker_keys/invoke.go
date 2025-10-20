package graphBetaRotateBitLockerKeys

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type rotateResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *RotateBitLockerKeysAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RotateBitLockerKeysActionModel

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
	tflog.Debug(ctx, fmt.Sprintf("Rotating BitLocker keys for %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting BitLocker key rotation for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Rotate keys concurrently with error collection
	results := make(chan rotateResult, totalDevices)
	var wg sync.WaitGroup

	// Rotate keys on managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.rotateManagedDevice(ctx, id)
			results <- rotateResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Rotate keys on co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.rotateComanagedDevice(ctx, id)
			results <- rotateResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to rotate BitLocker keys on %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully rotated BitLocker keys on %s device %s",
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
				fmt.Sprintf("Successfully rotated BitLocker keys on %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that had keys rotated now have new BitLocker recovery passwords escrowed in Intune/Azure AD. "+
					"Previous recovery keys are invalidated and can no longer be used for recovery. "+
					"Failed devices may be offline, not Windows devices, not have BitLocker enabled, or may have configuration issues.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully rotated BitLocker keys on %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("BitLocker key rotation completed for %d device(s). "+
				"New recovery passwords have been generated and escrowed to Intune/Azure AD. "+
				"Previous BitLocker recovery keys are now invalid and cannot be used for device recovery. "+
				"New recovery keys are available in the Microsoft Intune admin center and Azure AD portal. "+
				"No user action or device restart is required.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *RotateBitLockerKeysAction) rotateManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Rotating BitLocker keys on managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RotateBitLockerKeys().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to rotate BitLocker keys on managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *RotateBitLockerKeysAction) rotateComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Rotating BitLocker keys on co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RotateBitLockerKeys().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to rotate BitLocker keys on co-managed device %s: %w", deviceID, err)
	}

	return nil
}
