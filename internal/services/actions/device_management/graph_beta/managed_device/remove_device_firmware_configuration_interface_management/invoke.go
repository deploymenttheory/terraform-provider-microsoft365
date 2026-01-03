package graphBetaRemoveDeviceFirmwareConfigurationInterfaceManagementManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type removeDFCIResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceActionModel

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
	tflog.Debug(ctx, fmt.Sprintf("Removing DFCI management from %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting DFCI management removal for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Remove DFCI management from devices concurrently with error collection
	results := make(chan removeDFCIResult, totalDevices)
	var wg sync.WaitGroup

	// Remove DFCI from managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.removeDFCIManagedDevice(ctx, id)
			results <- removeDFCIResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Remove DFCI from co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.removeDFCIComanagedDevice(ctx, id)
			results <- removeDFCIResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to remove DFCI management from %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully removed DFCI management from %s device %s",
				result.deviceType, result.deviceID))

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): DFCI management removed", result.deviceID, result.deviceType),
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
				fmt.Sprintf("Successfully removed DFCI management from %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices with DFCI removed will no longer support remote UEFI/BIOS configuration via Intune. "+
					"Standard MDM management continues.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed DFCI management from %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("DFCI removal complete: %d device(s) removed from firmware-level management.\n\n"+
				"Note: These devices will continue standard Intune MDM management but will no longer support "+
				"remote UEFI/BIOS configuration. Physical access may be required to re-enable DFCI.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) removeDFCIManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing DFCI management from managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RemoveDeviceFirmwareConfigurationInterfaceManagement().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to remove DFCI management from managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) removeDFCIComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing DFCI management from co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RemoveDeviceFirmwareConfigurationInterfaceManagement().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to remove DFCI management from co-managed device %s: %w", deviceID, err)
	}

	return nil
}
