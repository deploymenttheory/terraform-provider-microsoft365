package graphBetaInitiateMobileDeviceManagementKeyRecoveryManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type keyRecoveryResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data InitiateMobileDeviceManagementKeyRecoveryManagedDeviceActionModel

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
	tflog.Debug(ctx, fmt.Sprintf("Initiating MDM key recovery for %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting MDM key recovery and TPM attestation for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Initiate key recovery concurrently with error collection
	results := make(chan keyRecoveryResult, totalDevices)
	var wg sync.WaitGroup

	// Initiate key recovery on managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.initiateKeyRecoveryManagedDevice(ctx, id)
			results <- keyRecoveryResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Initiate key recovery on co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.initiateKeyRecoveryComanagedDevice(ctx, id)
			results <- keyRecoveryResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to initiate MDM key recovery for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully initiated MDM key recovery for %s device %s",
				result.deviceType, result.deviceID))

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): MDM key recovery and TPM attestation initiated", result.deviceID, result.deviceType),
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
				fmt.Sprintf("Successfully initiated MDM key recovery for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices with successful key recovery will have their BitLocker recovery keys escrowed to Azure AD "+
					"and TPM attestation completed.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully initiated MDM key recovery for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ MDM key recovery complete: %d device(s) have BitLocker recovery keys escrowed to Azure AD "+
				"and TPM attestation completed.\n\n"+
				"Recovery keys are now available in Azure AD for disaster recovery scenarios. "+
				"This action does not affect device operation or user access.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) initiateKeyRecoveryManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Initiating MDM key recovery for managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		InitiateMobileDeviceManagementKeyRecovery().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate MDM key recovery for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) initiateKeyRecoveryComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Initiating MDM key recovery for co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		InitiateMobileDeviceManagementKeyRecovery().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate MDM key recovery for co-managed device %s: %w", deviceID, err)
	}

	return nil
}

