package graphBetaMoveDevicesToOUManagedDevice

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validation"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *MoveDevicesToOUManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data MoveDevicesToOUManagedDeviceActionModel

	tflog.Debug(ctx, "Starting move devices to OU action", map[string]any{"action": ActionName})

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

	ouPath := data.OrganizationalUnitPath.ValueString()

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
	tflog.Debug(ctx, "Processing devices for OU move", map[string]any{
		"managed_devices":   len(managedDeviceIDs),
		"comanaged_devices": len(comanagedDeviceIDs),
		"total_devices":     totalDevices,
		"ou_path":           ouPath,
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

		validationResult, err := validateRequest(ctx, a.client, managedDeviceIDs, comanagedDeviceIDs)
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
			Error(validationResult.NonExistentComanagedDevices, "co-managed device", "do not exist or are not co-managed by Intune").
			Warning(validationResult.UnsupportedManagedDevices, "managed device", "are not Windows devices (only Windows devices can be moved to AD OUs)").
			Warning(validationResult.UnsupportedComanagedDevices, "co-managed device", "are not Windows devices (only Windows devices can be moved to AD OUs)").
			Warning(validationResult.NotHybridJoinedDevices, "device", "may not be hybrid Azure AD joined (cloud-only or workplace-joined devices cannot be moved to AD OUs)")

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting move operation for %d device(s) to OU: %s", totalDevices, ouPath),
	})

	// Track successes and failures
	successCount := 0
	var moveErrors []error

	// Move managed devices (batch operation)
	if len(managedDeviceIDs) > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Moving %d managed device(s) to OU...", len(managedDeviceIDs)),
		})

		err := a.moveManagedDevicesToOU(ctx, managedDeviceIDs, ouPath)
		if err != nil {
			tflog.Error(ctx, "Failed to move managed devices to OU", map[string]any{"error": err.Error()})
			moveErrors = append(moveErrors, fmt.Errorf("managed devices: %w", err))
		} else {
			successCount += len(managedDeviceIDs)
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Successfully moved %d managed device(s) to OU: %s", len(managedDeviceIDs), ouPath),
			})
		}
	}

	// Move co-managed devices (batch operation)
	if len(comanagedDeviceIDs) > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Moving %d co-managed device(s) to OU...", len(comanagedDeviceIDs)),
		})

		err := a.moveComanagedDevicesToOU(ctx, comanagedDeviceIDs, ouPath)
		if err != nil {
			tflog.Error(ctx, "Failed to move co-managed devices to OU", map[string]any{"error": err.Error()})
			moveErrors = append(moveErrors, fmt.Errorf("co-managed devices: %w", err))
		} else {
			successCount += len(comanagedDeviceIDs)
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Successfully moved %d co-managed device(s) to OU: %s", len(comanagedDeviceIDs), ouPath),
			})
		}
	}

	// Handle results
	if len(moveErrors) > 0 {
		if successCount > 0 {
			// Partial success
			if ignorePartialFailures {
				resp.Diagnostics.AddWarning(
					"Partial Success (Failures Ignored)",
					fmt.Sprintf("Successfully moved %d of %d devices to OU '%s'. Errors: %v\n\n"+
						"Moved devices will appear in the new OU after the next Azure AD Connect sync cycle.",
						successCount, totalDevices, ouPath, moveErrors),
				)
			} else {
				resp.Diagnostics.AddError(
					"Partial Failure",
					fmt.Sprintf("Successfully moved %d of %d devices to OU '%s', but some operations failed. Errors: %v\n\n"+
						"Moved devices will appear in the new OU after the next Azure AD Connect sync cycle.",
						successCount, totalDevices, ouPath, moveErrors),
				)
				return
			}
		} else {
			// Complete failure
			if len(moveErrors) == 1 {
				errors.HandleKiotaGraphError(ctx, moveErrors[0], resp, "Action", a.WritePermissions)
			} else {
				resp.Diagnostics.AddError(
					"Failed to Move Devices to OU",
					fmt.Sprintf("Failed to move devices to OU '%s'. Multiple errors occurred: %v", ouPath, moveErrors),
				)
			}
			return
		}
	}

	if successCount > 0 {
		tflog.Info(ctx, "Move devices to OU action completed", map[string]any{
			"success_count": successCount,
			"total_devices": totalDevices,
			"ou_path":       ouPath,
		})

		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ Move to OU complete: %d device(s) moved to '%s'\n\n"+
				"The devices will appear in the new Organizational Unit after the next Azure AD Connect sync cycle. "+
				"This typically occurs within 30 minutes but may vary based on your sync schedule.",
				successCount, ouPath),
		})
	}

	tflog.Debug(ctx, "Finished move devices to OU action")
}

func (a *MoveDevicesToOUManagedDeviceAction) moveManagedDevicesToOU(ctx context.Context, deviceIDs []string, ouPath string) error {
	tflog.Debug(ctx, "Moving managed devices to OU", map[string]any{
		"device_count": len(deviceIDs),
		"ou_path":      ouPath,
	})

	requestBody := constructManagedDevicesRequest(ctx, deviceIDs, ouPath)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		MoveDevicesToOU().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}

func (a *MoveDevicesToOUManagedDeviceAction) moveComanagedDevicesToOU(ctx context.Context, deviceIDs []string, ouPath string) error {
	tflog.Debug(ctx, "Moving co-managed devices to OU", map[string]any{
		"device_count": len(deviceIDs),
		"ou_path":      ouPath,
	})

	requestBody := constructComanagedDevicesRequest(ctx, deviceIDs, ouPath)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		MoveDevicesToOU().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
