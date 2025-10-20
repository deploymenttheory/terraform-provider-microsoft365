package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data CreateDeviceLogCollectionRequestManagedDeviceActionModel

	tflog.Debug(ctx, "Starting device log collection action", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, "Processing devices for log collection", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
		"total_devices":     totalDevices,
	})

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting log collection requests for %d devices (%d managed, %d co-managed)",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	successCount := 0
	var failedDevices []string
	var lastError error

	// Process managed devices
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		templateType := "predefined"
		if !device.TemplateType.IsNull() && !device.TemplateType.IsUnknown() {
			templateType = device.TemplateType.ValueString()
		}

		err := a.createLogCollectionManagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (managed)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to create log collection request for managed device", map[string]any{
				"device_id":     deviceID,
				"template_type": templateType,
				"error":         err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully initiated log collection for managed device", map[string]any{
				"device_id":     deviceID,
				"template_type": templateType,
			})
		}

		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	// Process co-managed devices
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		templateType := "predefined"
		if !device.TemplateType.IsNull() && !device.TemplateType.IsUnknown() {
			templateType = device.TemplateType.ValueString()
		}

		err := a.createLogCollectionComanagedDevice(ctx, device)
		if err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (co-managed)", deviceID))
			lastError = err
			tflog.Error(ctx, "Failed to create log collection request for co-managed device", map[string]any{
				"device_id":     deviceID,
				"template_type": templateType,
				"error":         err.Error(),
			})
		} else {
			successCount++
			tflog.Debug(ctx, "Successfully initiated log collection for co-managed device", map[string]any{
				"device_id":     deviceID,
				"template_type": templateType,
			})
		}

		processed := successCount + len(failedDevices)
		progress := float64(processed) / float64(totalDevices) * 100
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Processed %d of %d devices (%.0f%% complete)",
				processed, totalDevices, progress),
		})
	}

	a.reportResults(ctx, resp, successCount, totalDevices, failedDevices, lastError)

	tflog.Debug(ctx, "Completed device log collection action", map[string]any{
		"action":        ActionName,
		"success_count": successCount,
		"failed_count":  len(failedDevices),
		"total_devices": totalDevices,
	})
}

// createLogCollectionManagedDevice performs atomic log collection for a managed device
func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) createLogCollectionManagedDevice(ctx context.Context, device ManagedDeviceLogCollection) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Creating log collection request for managed device", map[string]any{"device_id": deviceID})

	requestBody := constructManagedDeviceRequest(ctx, device)

	_, err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		CreateDeviceLogCollectionRequest().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated log collection for managed device", map[string]any{"device_id": deviceID})
	return nil
}

// createLogCollectionComanagedDevice performs atomic log collection for a co-managed device
func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) createLogCollectionComanagedDevice(ctx context.Context, device ComanagedDeviceLogCollection) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, "Creating log collection request for co-managed device", map[string]any{"device_id": deviceID})

	requestBody := constructComanagedDeviceRequest(ctx, device)

	_, err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		CreateDeviceLogCollectionRequest().
		Post(ctx, requestBody, nil)

	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully initiated log collection for co-managed device", map[string]any{"device_id": deviceID})
	return nil
}

// reportResults handles final result reporting according to ADR-001 principles
func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) reportResults(ctx context.Context, resp *action.InvokeResponse, successCount, totalDevices int, failedDevices []string, lastError error) {
	if len(failedDevices) > 0 {
		if successCount > 0 {
			// Partial success
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("Partial success: %d of %d devices succeeded. Failed devices: %v",
					successCount, totalDevices, failedDevices),
			})
			resp.Diagnostics.AddWarning(
				"Partial Success",
				fmt.Sprintf("Log collection requests created for %d of %d devices. Failed devices: %s. "+
					"Devices that received the command will collect logs and upload them to Intune.",
					successCount, totalDevices, strings.Join(failedDevices, ", ")),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	} else {
		// Full success
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Log collection requests initiated successfully for all %d devices. "+
				"Logs will be available in the Intune portal after collection completes.",
				successCount),
		})
	}
}
