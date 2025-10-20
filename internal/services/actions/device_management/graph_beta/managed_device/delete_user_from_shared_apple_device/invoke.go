package graphBetaDeleteUserFromSharedAppleDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *DeleteUserFromSharedAppleDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DeleteUserFromSharedAppleDeviceActionModel

	tflog.Debug(ctx, "Starting delete user from shared Apple device action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalOperations := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing device-user pairs for user deletion", map[string]any{
		"total_operations": totalOperations,
	})

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting deletion of %d user(s) from Shared iPad device(s)",
			totalOperations),
	})

	successCount := 0
	var failedOperations []string
	var lastError error

	// Process managed device-user pairs sequentially
	for _, deviceUser := range data.ManagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		upn := deviceUser.UserPrincipalName.ValueString()

		err := a.deleteUserFromManagedDevice(ctx, deviceID, upn)
		if err != nil {
			failedOperations = append(failedOperations, fmt.Sprintf("%s (Managed, User: %s)", deviceID, upn))
			lastError = err
			tflog.Error(ctx, "Failed to delete user from managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
				"error":               err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully deleted user from managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
			})
		}

		processed := successCount + len(failedOperations)
		progress := float64(processed) / float64(totalOperations) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d operations (%.0f%% complete)",
				processed, totalOperations, progress),
		})
	}

	// Process co-managed device-user pairs sequentially
	for _, deviceUser := range data.ComanagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		upn := deviceUser.UserPrincipalName.ValueString()

		err := a.deleteUserFromComanagedDevice(ctx, deviceID, upn)
		if err != nil {
			failedOperations = append(failedOperations, fmt.Sprintf("%s (Co-Managed, User: %s)", deviceID, upn))
			lastError = err
			tflog.Error(ctx, "Failed to delete user from co-managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
				"error":               err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully deleted user from co-managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
			})
		}

		processed := successCount + len(failedOperations)
		progress := float64(processed) / float64(totalOperations) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d operations (%.0f%% complete)",
				processed, totalOperations, progress),
		})
	}

	a.reportResults(ctx, resp, successCount, totalOperations, failedOperations, lastError)

	tflog.Debug(ctx, "Completed delete user from shared Apple device action", map[string]any{
		"action":           ActionName,
		"success_count":    successCount,
		"failed_count":     len(failedOperations),
		"total_operations": totalOperations,
	})
}

// deleteUserFromManagedDevice performs atomic user deletion from a managed Shared iPad device
func (a *DeleteUserFromSharedAppleDeviceAction) deleteUserFromManagedDevice(ctx context.Context, deviceID string, userPrincipalName string) error {
	tflog.Debug(ctx, "Deleting user from managed Shared iPad device", map[string]any{
		"device_id":           deviceID,
		"user_principal_name": userPrincipalName,
	})

	requestBody := constructManagedDeviceRequest(ctx, userPrincipalName)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		DeleteUserFromSharedAppleDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully deleted user from managed Shared iPad device", map[string]any{
		"device_id":           deviceID,
		"user_principal_name": userPrincipalName,
	})
	return nil
}

// deleteUserFromComanagedDevice performs atomic user deletion from a co-managed Shared iPad device
func (a *DeleteUserFromSharedAppleDeviceAction) deleteUserFromComanagedDevice(ctx context.Context, deviceID string, userPrincipalName string) error {
	tflog.Debug(ctx, "Deleting user from co-managed Shared iPad device", map[string]any{
		"device_id":           deviceID,
		"user_principal_name": userPrincipalName,
	})

	requestBody := constructComanagedDeviceRequest(ctx, userPrincipalName)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		DeleteUserFromSharedAppleDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully deleted user from co-managed Shared iPad device", map[string]any{
		"device_id":           deviceID,
		"user_principal_name": userPrincipalName,
	})
	return nil
}

// reportResults handles final result reporting according to ADR-001 principles
func (a *DeleteUserFromSharedAppleDeviceAction) reportResults(ctx context.Context, resp *action.InvokeResponse, successCount, totalOperations int, failedOperations []string, lastError error) {
	if len(failedOperations) > 0 {
		if successCount > 0 {
			// Partial success
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("Partial success: %d of %d operations succeeded. Failed operations: %v",
					successCount, totalOperations, failedOperations),
			})
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully deleted users from %d of %d operations. Failed operations: %s. "+
					"Successfully deleted users have been permanently removed from their respective Shared iPads, and all their cached data has been deleted. "+
					"Failed operations may be due to the device not being in Shared iPad mode, the user not existing on the device, or network issues.",
					successCount, totalOperations, strings.Join(failedOperations, ", ")),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	} else {
		// Full success
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Deletion complete: %d user(s) have been permanently removed from their Shared iPad device(s). "+
				"All cached user data has been deleted, and the freed storage space is now available for other users.",
				successCount),
		})
	}
}
