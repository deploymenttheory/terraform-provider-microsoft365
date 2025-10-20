package graphBetaLogoutSharedAppleDeviceActiveUser

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type logoutResult struct {
	deviceID string
	err      error
}

func (a *LogoutSharedAppleDeviceActiveUserAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data LogoutSharedAppleDeviceActiveUserActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, fmt.Sprintf("Performing action %s for %d Shared iPad device(s)", ActionName, totalDevices))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting logout of active users from %d Shared iPad device(s)...", totalDevices),
	})

	// Logout active users from Shared iPad devices concurrently with error collection
	results := make(chan logoutResult, totalDevices)
	var wg sync.WaitGroup

	for _, deviceID := range deviceIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := a.logoutActiveUser(ctx, id)
			results <- logoutResult{deviceID: id, err: err}
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
			failedDevices = append(failedDevices, result.deviceID)
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to logout active user from Shared iPad %s: %v", result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully logged out active user from Shared iPad %s", result.deviceID))
		}

		// Send progress update
		progress := float64(successCount+len(failedDevices)) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)", successCount+len(failedDevices), totalDevices, progress),
		})
	}

	// Report results
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully logged out active users from %d of %d Shared iPad devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices with successful logout have returned to the login screen. "+
					"Failed devices may not be in Shared iPad mode, may not have an active user, or may be offline.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully logged out active users from %d Shared iPad device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Logout complete: %d Shared iPad device(s) have returned to the login screen. "+
				"Devices are ready for the next user to log in. "+
				"User data remains cached on the devices.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *LogoutSharedAppleDeviceActiveUserAction) logoutActiveUser(ctx context.Context, deviceID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Logging out active user from Shared iPad with ID: %s", deviceID))

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		LogoutSharedAppleDeviceActiveUser().
		Post(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to logout active user from Shared iPad %s: %w", deviceID, err)
	}

	return nil
}
