package graphBetaRotateLocalAdminPasswordManagedDevice

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type passwordRotationResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RotateLocalAdminPasswordManagedDeviceActionModel

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
	tflog.Debug(ctx, fmt.Sprintf("Rotating local admin password for %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting local administrator password rotation for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Rotate passwords concurrently with error collection
	results := make(chan passwordRotationResult, totalDevices)
	var wg sync.WaitGroup

	// Rotate passwords on managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.rotatePasswordManagedDevice(ctx, id)
			results <- passwordRotationResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Rotate passwords on co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.rotatePasswordComanagedDevice(ctx, id)
			results <- passwordRotationResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to rotate local admin password for %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully initiated local admin password rotation for %s device %s",
				result.deviceType, result.deviceID))

			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Device %s (%s): Local admin password rotation initiated", result.deviceID, result.deviceType),
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
				fmt.Sprintf("Successfully rotated local administrator password for %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices with successful rotation will have new passwords stored in Azure AD/Intune. "+
					"Authorized administrators can retrieve the new passwords.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully rotated local admin password for %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ Password rotation complete: %d device(s) have new local administrator passwords.\n\n"+
				"Password details:\n"+
				"- New passwords automatically generated (complex, random)\n"+
				"- Passwords stored securely in Azure AD or Intune\n"+
				"- Previous passwords no longer valid\n"+
				"- Authorized administrators can retrieve new passwords via Azure Portal or Graph API\n"+
				"- Password retrieval requires appropriate permissions and auditing is enabled",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) rotatePasswordManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Rotating local admin password for managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RotateLocalAdminPassword().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to rotate local admin password for managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) rotatePasswordComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Rotating local admin password for co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RotateLocalAdminPassword().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to rotate local admin password for co-managed device %s: %w", deviceID, err)
	}

	return nil
}

