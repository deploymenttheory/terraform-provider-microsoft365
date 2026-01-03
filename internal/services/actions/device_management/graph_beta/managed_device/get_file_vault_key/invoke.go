package graphBetaGetFileVaultKeyManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type getKeyResult struct {
	deviceID    string
	deviceType  string // "managed" or "comanaged"
	recoveryKey string
	err         error
}

func (a *GetFileVaultKeyManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data GetFileVaultKeyManagedDeviceActionModel

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
	tflog.Debug(ctx, fmt.Sprintf("Retrieving FileVault keys for %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting FileVault key retrieval for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Retrieve keys concurrently with error collection
	results := make(chan getKeyResult, totalDevices)
	var wg sync.WaitGroup

	// Retrieve keys from managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			key, err := a.getFileVaultKeyManagedDevice(ctx, id)
			results <- getKeyResult{deviceID: id, deviceType: "managed", recoveryKey: key, err: err}
		}(deviceID)
	}

	// Retrieve keys from co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			key, err := a.getFileVaultKeyComanagedDevice(ctx, id)
			results <- getKeyResult{deviceID: id, deviceType: "comanaged", recoveryKey: key, err: err}
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
	var retrievedKeys []string

	for result := range results {
		if result.err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s)", result.deviceID, result.deviceType))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to retrieve FileVault key for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved FileVault key for %s device %s",
				result.deviceType, result.deviceID))

			// Display the retrieved key in progress message
			keyMessage := fmt.Sprintf("✓ Device %s (%s): FileVault Recovery Key = %s",
				result.deviceID, result.deviceType, result.recoveryKey)
			retrievedKeys = append(retrievedKeys, keyMessage)

			resp.SendProgress(action.InvokeProgressEvent{
				Message: keyMessage,
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
				fmt.Sprintf("Successfully retrieved FileVault keys for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Retrieved keys have been displayed in the action output above.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved FileVault keys for %d device(s)", successCount))

	if successCount > 0 {
		// Display summary with all keys
		summaryMessage := fmt.Sprintf("\n========== FileVault Recovery Keys Retrieved ==========\n"+
			"Total Devices: %d\n"+
			"Successful: %d\n"+
			"Failed: %d\n\n"+
			"SECURITY WARNING: The following recovery keys grant full access to encrypted device data.\n"+
			"Store these keys securely and ensure compliance with your organization's security policies.\n\n"+
			"%s\n"+
			"========================================================",
			totalDevices, successCount, len(failedDevices),
			fmt.Sprintf("%v", retrievedKeys))

		resp.SendProgress(action.InvokeProgressEvent{
			Message: summaryMessage,
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *GetFileVaultKeyManagedDeviceAction) getFileVaultKeyManagedDevice(ctx context.Context, deviceID string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving FileVault key for managed device with ID: %s", deviceID))

	result, err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		GetFileVaultKey().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to retrieve FileVault key for managed device %s: %w", deviceID, err)
	}

	if result == nil || result.GetValue() == nil {
		return "", fmt.Errorf("no FileVault key returned for managed device %s", deviceID)
	}

	return *result.GetValue(), nil
}

func (a *GetFileVaultKeyManagedDeviceAction) getFileVaultKeyComanagedDevice(ctx context.Context, deviceID string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving FileVault key for co-managed device with ID: %s", deviceID))

	result, err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		GetFileVaultKey().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to retrieve FileVault key for co-managed device %s: %w", deviceID, err)
	}

	if result == nil || result.GetValue() == nil {
		return "", fmt.Errorf("no FileVault key returned for co-managed device %s", deviceID)
	}

	return *result.GetValue(), nil
}
