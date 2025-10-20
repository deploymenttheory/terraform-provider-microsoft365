package graphBetaMoveDevicesToOUManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *MoveDevicesToOUManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data MoveDevicesToOUManagedDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	tflog.Debug(ctx, fmt.Sprintf("Moving %d managed device(s) and %d co-managed device(s) to OU: %s",
		len(managedDeviceIDs), len(comanagedDeviceIDs), ouPath))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting move operation for %d device(s) to OU: %s",
			totalDevices, ouPath),
	})

	// Track successes and failures
	successCount := 0
	var errors []error

	// Move managed devices (collection-level operation)
	if len(managedDeviceIDs) > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Moving %d managed device(s) to OU...", len(managedDeviceIDs)),
		})

		err := a.moveManagedDevicesToOU(ctx, managedDeviceIDs, ouPath)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to move managed devices to OU: %v", err))
			errors = append(errors, fmt.Errorf("managed devices: %w", err))
		} else {
			successCount += len(managedDeviceIDs)
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Successfully moved %d managed device(s) to OU: %s",
					len(managedDeviceIDs), ouPath),
			})
		}
	}

	// Move co-managed devices (collection-level operation)
	if len(comanagedDeviceIDs) > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Moving %d co-managed device(s) to OU...", len(comanagedDeviceIDs)),
		})

		err := a.moveComanagedDevicesToOU(ctx, comanagedDeviceIDs, ouPath)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to move co-managed devices to OU: %v", err))
			errors = append(errors, fmt.Errorf("co-managed devices: %w", err))
		} else {
			successCount += len(comanagedDeviceIDs)
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("✓ Successfully moved %d co-managed device(s) to OU: %s",
					len(comanagedDeviceIDs), ouPath),
			})
		}
	}

	// Report results
	if len(errors) > 0 {
		if successCount > 0 {
			// Partial success
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Successfully moved %d of %d devices to OU '%s'. Errors: %v\n\n"+
					"Moved devices will appear in the new OU after the next Azure AD Connect sync cycle.",
					successCount, totalDevices, ouPath, errors),
			)
		} else {
			// Complete failure
			if len(errors) == 1 {
				handleMoveDevicesToOUError(ctx, errors[0], resp, a.WritePermissions)
			} else {
				resp.Diagnostics.AddError(
					"Failed to Move Devices to OU",
					fmt.Sprintf("Failed to move devices to OU '%s'. Multiple errors occurred: %v", ouPath, errors),
				)
			}
			return
		}
	}

	if successCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Successfully moved %d device(s) to OU: %s", successCount, ouPath))
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("✓ Move to OU complete: %d device(s) moved to '%s'\n\n"+
				"The devices will appear in the new Organizational Unit after the next Azure AD Connect sync cycle. "+
				"This typically occurs within 30 minutes but may vary based on your sync schedule.",
				successCount, ouPath),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *MoveDevicesToOUManagedDeviceAction) moveManagedDevicesToOU(ctx context.Context, deviceIDs []string, ouPath string) error {
	tflog.Debug(ctx, fmt.Sprintf("Moving %d managed devices to OU: %s", len(deviceIDs), ouPath))

	requestBody := constructManagedDevicesRequest(ctx, deviceIDs, ouPath)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		MoveDevicesToOU().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to move %d managed device(s) to OU '%s': %w", len(deviceIDs), ouPath, err)
	}

	return nil
}

func (a *MoveDevicesToOUManagedDeviceAction) moveComanagedDevicesToOU(ctx context.Context, deviceIDs []string, ouPath string) error {
	tflog.Debug(ctx, fmt.Sprintf("Moving %d co-managed devices to OU: %s", len(deviceIDs), ouPath))

	requestBody := constructComanagedDevicesRequest(ctx, deviceIDs, ouPath)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		MoveDevicesToOU().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to move %d co-managed device(s) to OU '%s': %w", len(deviceIDs), ouPath, err)
	}

	return nil
}

func handleMoveDevicesToOUError(ctx context.Context, err error, resp *action.InvokeResponse, permissions []string) {
	errors.HandleKiotaGraphError(ctx, err, resp, "Action", permissions)
}
