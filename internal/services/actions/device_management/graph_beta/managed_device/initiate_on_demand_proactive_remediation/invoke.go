package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

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

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data InitiateOnDemandProactiveRemediationManagedDeviceActionModel

	tflog.Debug(ctx, "Starting on-demand proactive remediation action", map[string]any{"action": ActionName})

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

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for proactive remediation", map[string]any{
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not Windows devices").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Windows devices")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("on-demand proactive remediation",
			fmt.Sprintf("%d managed, %d co-managed", len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		scriptPolicyID := device.ScriptPolicyID.ValueString()

		err := a.initiateRemediationManagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Managed").Failed(fmt.Sprintf("script policy %s: %s", scriptPolicyID, err.Error()))
			tflog.Error(ctx, "Failed to initiate remediation for managed device", map[string]any{
				"device_id":        deviceID,
				"script_policy_id": scriptPolicyID,
				"error":            err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Managed").Succeeded(fmt.Sprintf("proactive remediation initiated for script policy %s", scriptPolicyID))
			tflog.Info(ctx, "Successfully initiated remediation for managed device", map[string]any{
				"device_id":        deviceID,
				"script_policy_id": scriptPolicyID,
			})
		}
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		scriptPolicyID := device.ScriptPolicyID.ValueString()

		err := a.initiateRemediationComanagedDevice(ctx, device)
		if err != nil {
			progressTracker.Device(deviceID, "Co-managed").Failed(fmt.Sprintf("script policy %s: %s", scriptPolicyID, err.Error()))
			tflog.Error(ctx, "Failed to initiate remediation for co-managed device", map[string]any{
				"device_id":        deviceID,
				"script_policy_id": scriptPolicyID,
				"error":            err.Error(),
			})
		} else {
			progressTracker.Device(deviceID, "Co-managed").Succeeded(fmt.Sprintf("proactive remediation initiated for script policy %s", scriptPolicyID))
			tflog.Info(ctx, "Successfully initiated remediation for co-managed device", map[string]any{
				"device_id":        deviceID,
				"script_policy_id": scriptPolicyID,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("on-demand proactive remediation")
			tflog.Warn(ctx, "Proactive remediation completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Proactive Remediation Failed", "initiate on-demand proactive remediation on devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("initiated on-demand proactive remediation on all devices")
	}

	tflog.Info(ctx, "On-demand proactive remediation action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) initiateRemediationManagedDevice(ctx context.Context, device ManagedDeviceProactiveRemediation) error {
	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(device.DeviceID.ValueString()).
		InitiateOnDemandProactiveRemediation().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) initiateRemediationComanagedDevice(ctx context.Context, device ComanagedDeviceProactiveRemediation) error {
	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(device.DeviceID.ValueString()).
		InitiateOnDemandProactiveRemediation().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
