package graphBetaRevokeAppleVppLicenses

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type revokeResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	err        error
}

func (a *RevokeAppleVppLicensesAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RevokeAppleVppLicensesActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)
	tflog.Debug(ctx, fmt.Sprintf("Revoking Apple VPP licenses for %d managed device(s) and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Apple VPP license revocation for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(managedDeviceIDs), len(comanagedDeviceIDs)),
	})

	// Revoke licenses concurrently with error collection
	results := make(chan revokeResult, totalDevices)
	var wg sync.WaitGroup

	// Revoke from managed devices
	for _, deviceID := range managedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.revokeManagedDevice(ctx, id)
			results <- revokeResult{deviceID: id, deviceType: "managed", err: err}
		}(deviceID)
	}

	// Revoke from co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.revokeComanagedDevice(ctx, id)
			results <- revokeResult{deviceID: id, deviceType: "comanaged", err: err}
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
			tflog.Error(ctx, fmt.Sprintf("Failed to revoke Apple VPP licenses from %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully revoked Apple VPP licenses from %s device %s",
				result.deviceType, result.deviceID))
		}

		// Send progress update
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
				fmt.Sprintf("Successfully revoked Apple VPP licenses from %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that had licenses revoked will have their VPP app licenses returned to the available pool. "+
					"These licenses can now be reassigned to other devices. VPP apps may be removed from devices. "+
					"Failed devices may be offline, not iOS/iPadOS devices, or may not have any VPP licenses assigned.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully revoked Apple VPP licenses from %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Apple VPP license revocation completed for %d device(s). "+
				"All VPP-purchased app licenses have been revoked from these devices and returned to the available license pool. "+
				"Licenses are now available for reassignment to other devices or users. "+
				"VPP apps may be removed from the devices. Check the Apple Business Manager portal for updated license counts.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *RevokeAppleVppLicensesAction) revokeManagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Revoking Apple VPP licenses from managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		RevokeAppleVppLicenses().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to revoke Apple VPP licenses from managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *RevokeAppleVppLicensesAction) revokeComanagedDevice(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Revoking Apple VPP licenses from co-managed device with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		RevokeAppleVppLicenses().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to revoke Apple VPP licenses from co-managed device %s: %w", deviceID, err)
	}

	return nil
}
