package graphBetaDeleteUserFromSharedAppleDevice

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

func (a *DeleteUserFromSharedAppleDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DeleteUserFromSharedAppleDeviceActionModel

	tflog.Debug(ctx, "Starting delete user from shared Apple device action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	totalOperations := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing device-user pairs for user deletion", map[string]any{
		"total_operations": totalOperations,
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
			Error(validationResult.UnsupportedManagedDevices, "managed device", "are not Shared iPad devices (this action only supports Shared iPad devices)").
			Error(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Shared iPad devices (this action only supports Shared iPad devices)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Create progress tracker and send initial message
	progressTracker := progress.For(resp).WithTotalDevices(totalOperations).
		Starting("user deletion from Shared iPad devices",
			fmt.Sprintf("%d managed, %d co-managed", len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Process managed device-user pairs
	for _, deviceUser := range data.ManagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		upn := deviceUser.UserPrincipalName.ValueString()
		err := a.deleteUserFromManagedDevice(ctx, deviceID, upn)
		if err != nil {
			progressTracker.Device(fmt.Sprintf("%s (User: %s)", deviceID, upn), "Managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to delete user from managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
				"error":               err.Error(),
			})
		} else {
			progressTracker.Device(fmt.Sprintf("%s (User: %s)", deviceID, upn), "Managed").Succeeded("user deleted successfully")
			tflog.Info(ctx, "Successfully deleted user from managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
			})
		}
	}

	// Process co-managed device-user pairs
	for _, deviceUser := range data.ComanagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()
		upn := deviceUser.UserPrincipalName.ValueString()
		err := a.deleteUserFromComanagedDevice(ctx, deviceID, upn)
		if err != nil {
			progressTracker.Device(fmt.Sprintf("%s (User: %s)", deviceID, upn), "Co-managed").Failed(err.Error())
			tflog.Error(ctx, "Failed to delete user from co-managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
				"error":               err.Error(),
			})
		} else {
			progressTracker.Device(fmt.Sprintf("%s (User: %s)", deviceID, upn), "Co-managed").Succeeded("user deleted successfully")
			tflog.Info(ctx, "Successfully deleted user from co-managed Shared iPad device", map[string]any{
				"device_id":           deviceID,
				"user_principal_name": upn,
			})
		}
	}

	// Handle results
	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("user deletion from Shared iPad devices")
			tflog.Warn(ctx, "User deletion completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("User Deletion Failed", "delete users from Shared iPad devices")
			return
		}
	} else {
		progressTracker.CompletedSuccessfully("deleted all users from Shared iPad devices. All cached user data has been removed and storage space is now available")
	}

	tflog.Info(ctx, "Delete user from shared Apple device action completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_operations":         totalOperations,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

// deleteUserFromManagedDevice performs user deletion from a managed Shared iPad device
func (a *DeleteUserFromSharedAppleDeviceAction) deleteUserFromManagedDevice(ctx context.Context, deviceID string, userPrincipalName string) error {
	requestBody := constructManagedDeviceRequest(ctx, userPrincipalName)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		DeleteUserFromSharedAppleDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

// deleteUserFromComanagedDevice performs user deletion from a co-managed Shared iPad device
func (a *DeleteUserFromSharedAppleDeviceAction) deleteUserFromComanagedDevice(ctx context.Context, deviceID string, userPrincipalName string) error {
	requestBody := constructComanagedDeviceRequest(ctx, userPrincipalName)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		DeleteUserFromSharedAppleDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
