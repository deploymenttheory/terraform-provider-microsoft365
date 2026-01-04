package graphBetaActivateDeviceEsimManagedDevice

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

// Invoke runs the eSIM activation logic
func (a *ActivateDeviceEsimManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ActivateDeviceEsimManagedDeviceActionModelV2

	tflog.Debug(ctx, "Starting eSIM activation action V2", map[string]any{"action": ActionName})

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
	tflog.Debug(ctx, "Processing devices for eSIM activation", map[string]any{
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "have unsupported OS types (only iOS and iPadOS are supported)").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "have unsupported OS types (only iOS and iPadOS are supported)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("eSIM activation",
			fmt.Sprintf("%d managed, %d co-managed", len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.activateEsimManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to activate eSIM for managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Managed").Succeeded("eSIM activated successfully")
			tflog.Info(ctx, "Successfully activated eSIM for managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		err := a.activateEsimComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to activate eSIM for co-managed device", map[string]any{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Co-managed").Succeeded("eSIM activated successfully")
			tflog.Info(ctx, "Successfully activated eSIM for co-managed device", map[string]any{
				"device_id": deviceID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("eSIM activation")
			tflog.Warn(ctx, "eSIM activation completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("eSIM Activation Failed", "activate eSIM")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("activated eSIM")
	}

	tflog.Info(ctx, "eSIM activation action completed successfully", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"managed_devices":          len(data.ManagedDevices),
		"comanaged_devices":        len(data.ComanagedDevices),
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

// activateEsimManagedDevice activates eSIM on a managed device
func (a *ActivateDeviceEsimManagedDeviceAction) activateEsimManagedDevice(ctx context.Context, device ManagedDeviceActivateEsim) error {
	deviceID := device.DeviceID.ValueString()
	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		ActivateDeviceEsim().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

// activateEsimComanagedDevice activates eSIM on a co-managed device
func (a *ActivateDeviceEsimManagedDeviceAction) activateEsimComanagedDevice(ctx context.Context, device ComanagedDeviceActivateEsim) error {
	deviceID := device.DeviceID.ValueString()
	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		ActivateDeviceEsim().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
