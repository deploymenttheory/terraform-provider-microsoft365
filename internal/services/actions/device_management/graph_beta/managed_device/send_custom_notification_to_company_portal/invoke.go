package graphBetaSendCustomNotificationToCompanyPortal

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type notificationResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	title      string
	err        error
}

func (a *SendCustomNotificationToCompanyPortalAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SendCustomNotificationToCompanyPortalActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Sending custom notifications to %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting custom notification delivery to %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Send notifications concurrently with error collection
	results := make(chan notificationResult, totalDevices)
	var wg sync.WaitGroup

	// Send notifications to managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceNotification) {
			defer wg.Done()
			err := a.sendToManagedDevice(ctx, d)
			results <- notificationResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "managed",
				title:      d.NotificationTitle.ValueString(),
				err:        err,
			}
		}(device)
	}

	// Send notifications to co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceNotification) {
			defer wg.Done()
			err := a.sendToComanagedDevice(ctx, d)
			results <- notificationResult{
				deviceID:   d.DeviceID.ValueString(),
				deviceType: "comanaged",
				title:      d.NotificationTitle.ValueString(),
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
			tflog.Error(ctx, fmt.Sprintf("Failed to send notification to %s device %s: %v",
				result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			tflog.Debug(ctx, fmt.Sprintf("Successfully sent notification '%s' to %s device %s",
				result.title, result.deviceType, result.deviceID))
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
				fmt.Sprintf("Successfully sent custom notifications to %d of %d devices. Failed devices: %v. Last error: %v\n\n"+
					"Devices that received notifications will display them in the Company Portal app. "+
					"Failed devices may be offline, not have Company Portal installed, or the user may not be signed into Company Portal. "+
					"Users will see the notification the next time they open the Company Portal app.",
					successCount, totalDevices, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully sent custom notifications to %d device(s)", successCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Custom notification delivery completed for %d device(s). "+
				"Notifications have been sent to the Company Portal app on the target devices. "+
				"Users will see the notifications when they open or check the Company Portal app. "+
				"Notifications appear in the Company Portal notifications section and may also trigger push notifications "+
				"depending on the user's notification settings and device configuration.",
				successCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *SendCustomNotificationToCompanyPortalAction) sendToManagedDevice(ctx context.Context, device ManagedDeviceNotification) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Sending custom notification to managed device with ID: %s", deviceID))

	requestBody := constructManagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		SendCustomNotificationToCompanyPortal().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to send custom notification to managed device %s: %w", deviceID, err)
	}

	return nil
}

func (a *SendCustomNotificationToCompanyPortalAction) sendToComanagedDevice(ctx context.Context, device ComanagedDeviceNotification) error {
	deviceID := device.DeviceID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Sending custom notification to co-managed device with ID: %s", deviceID))

	requestBody := constructComanagedDeviceRequest(ctx, device)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		SendCustomNotificationToCompanyPortal().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to send custom notification to co-managed device %s: %w", deviceID, err)
	}

	return nil
}
