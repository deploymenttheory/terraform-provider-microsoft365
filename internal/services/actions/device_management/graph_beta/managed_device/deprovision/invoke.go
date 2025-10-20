package graphBetaDeprovisionManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type deprovisionResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	reason     string
	err        error
}

func (a *DeprovisionManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DeprovisionManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Deprovisioning %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting deprovision for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Deprovision devices concurrently with error collection
	results := make(chan deprovisionResult, totalDevices)
	var wg sync.WaitGroup

	// Deprovision managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceDeprovision) {
			defer wg.Done()
			err := a.deprovisionManagedDevice(ctx, d)
			results <- deprovisionResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "managed",
				reason:     d.DeprovisionReason.ValueString(),
				err:        err,
			}
		}(device)
	}

	// Deprovision co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceDeprovision) {
			defer wg.Done()
			err := a.deprovisionComanagedDevice(ctx, d)
			results <- deprovisionResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "comanaged",
				reason:     d.DeprovisionReason.ValueString(),
				err:        err,
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
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s)", result.deviceID, result.deviceType))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to deprovision %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully deprovisioned %s device %s (reason: %s)",
				result.deviceType, result.deviceID, result.reason))

			// Display success message
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Deprovisioned - Reason: %s",
					result.deviceID, result.deviceType, result.reason),
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
				fmt.Sprintf("Successfully deprovisioned %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Deprovisioned devices will have management policies and profiles removed.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deprovisioned %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Deprovision complete: %d device(s) have been deprovisioned from management.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *DeprovisionManagedDeviceAction) deprovisionManagedDevice(ctx context.Context, device ManagedDeviceDeprovision) error {
	deviceID := device.DeviceID.ValueString()
	reason := device.DeprovisionReason.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deprovisioning managed device with ID: %s (reason: %s)", deviceID, reason))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		Deprovision().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to deprovision managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *DeprovisionManagedDeviceAction) deprovisionComanagedDevice(ctx context.Context, device ComanagedDeviceDeprovision) error {
	deviceID := device.DeviceID.ValueString()
	reason := device.DeprovisionReason.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deprovisioning co-managed device with ID: %s (reason: %s)", deviceID, reason))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		Deprovision().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to deprovision co-managed device %s: %w", deviceID, err)
	}

	return nil
}
