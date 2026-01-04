package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/progress"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validation"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data CreateDeviceLogCollectionRequestManagedDeviceActionModel

	tflog.Debug(ctx, "Starting device log collection action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for log collection", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
		"total_devices":     totalDevices,
	})

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

		validationResult, err := validateRequest(ctx, a.client, data.ManagedDevices, data.ComanagedDevices)
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not Windows devices or do not support log collection (requires Windows 10 version 1709 or later, or Windows 11)").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Windows devices or do not support log collection (requires Windows 10 version 1709 or later, or Windows 11)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("device log collection",
			fmt.Sprintf("%d managed, %d co-managed", len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.createLogCollectionManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to create log collection request for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Managed").Succeeded("log collection request initiated successfully")
			tflog.Info(ctx, "Successfully initiated log collection for managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.createLogCollectionComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to create log collection request for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Co-managed").Succeeded("log collection request initiated successfully")
			tflog.Info(ctx, "Successfully initiated log collection for co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("device log collection")
			tflog.Warn(ctx, "Device log collection completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Device Log Collection Failed", "create log collection requests")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("initiated log collection requests for all devices. Logs will be available in the Intune portal after collection completes")
	}

	tflog.Info(ctx, "Device log collection action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

// createLogCollectionManagedDevice performs log collection request for a managed device
func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) createLogCollectionManagedDevice(ctx context.Context, device ManagedDeviceLogCollection) error {
	deviceID := device.DeviceID.ValueString()
	requestBody := constructManagedDeviceRequest(ctx, device)

	_, err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		CreateDeviceLogCollectionRequest().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

// createLogCollectionComanagedDevice performs log collection request for a co-managed device
func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) createLogCollectionComanagedDevice(ctx context.Context, device ComanagedDeviceLogCollection) error {
	deviceID := device.DeviceID.ValueString()
	requestBody := constructComanagedDeviceRequest(ctx, device)

	_, err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		CreateDeviceLogCollectionRequest().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
