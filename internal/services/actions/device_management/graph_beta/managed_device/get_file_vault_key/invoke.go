package graphBetaGetFileVaultKeyManagedDevice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/progress"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validation"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *GetFileVaultKeyManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data GetFileVaultKeyManagedDeviceActionModel

	tflog.Debug(ctx, "Starting FileVault key retrieval action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle timeout
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

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
	tflog.Debug(ctx, "Processing devices for FileVault key retrieval", map[string]any{
		"total_devices": totalDevices,
	})

	// Get ignore_partial_failures setting
	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	// Get validate_device_exists setting (default: true)
	validateDeviceExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateDeviceExists = data.ValidateDeviceExists.ValueBool()
	}

	// Perform API validation of devices if enabled
	if validateDeviceExists {
		tflog.Debug(ctx, "Performing device validation via API")

		validationResult, err := validateRequest(ctx, a.client, managedDeviceIDs, comanagedDeviceIDs)
		if err != nil {
			tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError(
				"Device Validation Failed",
				fmt.Sprintf("Failed to validate devices: %s", err.Error()),
			)
			return
		}

		// Report validation results
		results := validation.NewResults().
			Error(validationResult.NonExistentManagedDevices, "managed device", "do not exist or are not managed by Intune").
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not managed by Intune").
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not macOS devices or do not have FileVault enabled").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not macOS devices or do not have FileVault enabled")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Security warning
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "⚠️  SECURITY WARNING: FileVault recovery keys will be displayed in plain text. " +
			"These keys grant full access to encrypted device data. Ensure proper security controls are in place.",
	})

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("FileVault key retrieval",
			fmt.Sprintf("%d managed, %d co-managed", len(managedDeviceIDs), len(comanagedDeviceIDs)))

	// Track retrieved keys for summary
	var retrievedKeys []string

	// Process managed devices
	for _, deviceID := range managedDeviceIDs {
		key, err := a.getFileVaultKeyManagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to retrieve FileVault key for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			keyMessage := fmt.Sprintf("FileVault Recovery Key = %s", key)
			progressTracker.Device(deviceID, "Managed").Succeeded(keyMessage)
			retrievedKeys = append(retrievedKeys, fmt.Sprintf("Device %s (Managed): %s", deviceID, key))
			tflog.Info(ctx, "Successfully retrieved FileVault key for managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		key, err := a.getFileVaultKeyComanagedDevice(ctx, deviceID)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to retrieve FileVault key for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			keyMessage := fmt.Sprintf("FileVault Recovery Key = %s", key)
			progressTracker.Device(deviceID, "Co-managed").Succeeded(keyMessage)
			retrievedKeys = append(retrievedKeys, fmt.Sprintf("Device %s (Co-managed): %s", deviceID, key))
			tflog.Info(ctx, "Successfully retrieved FileVault key for co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("FileVault key retrieval")
			tflog.Warn(ctx, "FileVault key retrieval completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("FileVault Key Retrieval Failed", "retrieve FileVault keys from devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("retrieved all FileVault recovery keys")
	}

	// Display summary with all keys if any were retrieved
	if len(retrievedKeys) > 0 {
		summaryMessage := fmt.Sprintf("\n========== FileVault Recovery Keys Retrieved ==========\n"+
			"Total Devices: %d\n"+
			"Successful: %d\n"+
			"Failed: %d\n\n"+
			"⚠️  SECURITY WARNING: The following recovery keys grant full access to encrypted device data.\n"+
			"Store these keys securely and ensure compliance with your organization's security policies.\n\n"+
			"%s\n"+
			"========================================================",
			totalDevices, progressTracker.SuccessCount(), progressTracker.FailureCount(),
			strings.Join(retrievedKeys, "\n"))

		resp.SendProgress(action.InvokeProgressEvent{
			Message: summaryMessage,
		})
	}

	tflog.Info(ctx, "FileVault key retrieval action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *GetFileVaultKeyManagedDeviceAction) getFileVaultKeyManagedDevice(ctx context.Context, deviceID string) (string, error) {
	result, err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		GetFileVaultKey().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	if result == nil || result.GetValue() == nil {
		return "", fmt.Errorf("no FileVault key returned for device")
	}

	return *result.GetValue(), nil
}

func (a *GetFileVaultKeyManagedDeviceAction) getFileVaultKeyComanagedDevice(ctx context.Context, deviceID string) (string, error) {
	result, err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		GetFileVaultKey().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	if result == nil || result.GetValue() == nil {
		return "", fmt.Errorf("no FileVault key returned for device")
	}

	return *result.GetValue(), nil
}
