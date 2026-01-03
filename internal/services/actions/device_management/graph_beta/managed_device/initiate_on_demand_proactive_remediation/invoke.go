package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type remediationResult struct {
	deviceID       string
	scriptPolicyID string
	deviceType     string // "managed" or "comanaged"
	err            error
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data InitiateOnDemandProactiveRemediationManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Initiating on-demand proactive remediation for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting on-demand proactive remediation for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Initiate remediation concurrently with error collection
	results := make(chan remediationResult, totalDevices)
	var wg sync.WaitGroup

	// Initiate remediation on managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceProactiveRemediation) {
			defer wg.Done()
			err := a.initiateRemediationManagedDevice(ctx, d)
			results <- remediationResult{
				deviceID:       d.DeviceID.ValueString(),
				scriptPolicyID: d.ScriptPolicyID.ValueString(),
				deviceType:     "managed",
				err:            err,
			}
		}(device)
	}

	// Initiate remediation on co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceProactiveRemediation) {
			defer wg.Done()
			err := a.initiateRemediationComanagedDevice(ctx, d)
			results <- remediationResult{
				deviceID:       d.DeviceID.ValueString(),
				scriptPolicyID: d.ScriptPolicyID.ValueString(),
				deviceType:     "comanaged",
				err:            err,
			}
		}(device)
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
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s, script: %s)",
				result.deviceID, result.deviceType, result.scriptPolicyID))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to initiate remediation for %s device %s with script %s: %v",
				result.deviceType, result.deviceID, result.scriptPolicyID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully initiated remediation for %s device %s with script %s",
				result.deviceType, result.deviceID, result.scriptPolicyID))

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Remediation initiated (Script: %s)",
					result.deviceID, result.deviceType, result.scriptPolicyID),
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
				fmt.Sprintf("Successfully initiated on-demand proactive remediation for %d of %d devices. "+
					"Failed devices: %v. Last error: %v\n\n"+
					"Devices with successful initiation will execute the remediation script at next check-in. "+
					"Results will be available in Intune portal.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully initiated on-demand proactive remediation for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ On-demand proactive remediation initiated: %d device(s) will execute remediation scripts at next check-in.\n\n"+
				"Script execution details:\n"+
				"- Scripts run with SYSTEM privileges\n"+
				"- Execution may take several minutes\n"+
				"- Results available in Intune → Devices → Remediations\n"+
				"- Check device-specific logs for detailed execution status",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) initiateRemediationManagedDevice(ctx context.Context, device ManagedDeviceProactiveRemediation) error {
	deviceID := device.DeviceID.ValueString()
	scriptPolicyID := device.ScriptPolicyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Initiating on-demand proactive remediation for managed device %s with script policy %s",
		deviceID, scriptPolicyID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		InitiateOnDemandProactiveRemediation().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate on-demand proactive remediation for managed device %s with script policy %s: %w",
			deviceID, scriptPolicyID, err)
	}

	return nil
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) initiateRemediationComanagedDevice(ctx context.Context, device ComanagedDeviceProactiveRemediation) error {
	deviceID := device.DeviceID.ValueString()
	scriptPolicyID := device.ScriptPolicyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Initiating on-demand proactive remediation for co-managed device %s with script policy %s",
		deviceID, scriptPolicyID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		InitiateOnDemandProactiveRemediation().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate on-demand proactive remediation for co-managed device %s with script policy %s: %w",
			deviceID, scriptPolicyID, err)
	}

	return nil
}
