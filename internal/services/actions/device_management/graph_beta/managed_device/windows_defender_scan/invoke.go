package graphBetaWindowsDefenderScan

import (
	"context"
	"fmt"
	"sync"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type scanResult struct {
	deviceID   string
	deviceType string // "managed" or "comanaged"
	scanType   string // "quick" or "full"
	err        error
}

func (a *WindowsDefenderScanAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data WindowsDefenderScanActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	tflog.Debug(ctx, fmt.Sprintf("Performing Windows Defender scan for %d managed device(s) and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting Windows Defender scan for %d device(s) (%d managed, %d co-managed)...",
			totalDevices, len(data.ManagedDevices), len(data.ComanagedDevices)),
	})

	// Scan devices concurrently with error collection
	results := make(chan scanResult, totalDevices)
	var wg sync.WaitGroup

	// Scan managed devices
	for _, device := range data.ManagedDevices {
		wg.Add(1)
		go func(d ManagedDeviceScan) {
			defer wg.Done()
			deviceID := d.DeviceID.ValueString()
			quickScan := d.QuickScan.ValueBool()
			scanType := "full"
			if quickScan {
				scanType = "quick"
			}
			err := a.scanManagedDevice(ctx, deviceID, quickScan)
			results <- scanResult{deviceID: deviceID, deviceType: "managed", scanType: scanType, err: err}
		}(device)
	}

	// Scan co-managed devices
	for _, device := range data.ComanagedDevices {
		wg.Add(1)
		go func(d ComanagedDeviceScan) {
			defer wg.Done()
			deviceID := d.DeviceID.ValueString()
			quickScan := d.QuickScan.ValueBool()
			scanType := "full"
			if quickScan {
				scanType = "quick"
			}
			err := a.scanComanagedDevice(ctx, deviceID, quickScan)
			results <- scanResult{deviceID: deviceID, deviceType: "comanaged", scanType: scanType, err: err}
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
	quickScanCount := 0
	fullScanCount := 0

	for result := range results {
		if result.err != nil {
			failedDevices = append(failedDevices, fmt.Sprintf("%s (%s, %s scan)", result.deviceID, result.deviceType, result.scanType))
			lastError = result.err
			tflog.Error(ctx, fmt.Sprintf("Failed to initiate %s scan on %s device %s: %v",
				result.scanType, result.deviceType, result.deviceID, result.err))
		} else {
			successCount++
			if result.scanType == "quick" {
				quickScanCount++
			} else {
				fullScanCount++
			}
			tflog.Debug(ctx, fmt.Sprintf("Successfully initiated %s scan on %s device %s",
				result.scanType, result.deviceType, result.deviceID))
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
				fmt.Sprintf("Successfully initiated scans on %d of %d devices (%d quick scans, %d full scans). "+
					"Failed devices: %v. Last error: %v\n\n"+
					"Devices that received the scan command will begin scanning immediately if online. "+
					"Quick scans take 5-15 minutes, full scans take 30+ minutes to hours. "+
					"Failed devices may be offline, not Windows devices, or may not have Windows Defender enabled.",
					successCount, totalDevices, quickScanCount, fullScanCount, failedDevices, lastError),
			)
		} else {
			// Complete failure
			errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully initiated scans on %d device(s) (%d quick, %d full)",
		successCount, quickScanCount, fullScanCount))

	if successCount > 0 {
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("Windows Defender scan initiated on %d device(s): "+
				"%d quick scan(s) (5-15 minutes), %d full scan(s) (30+ minutes to hours). "+
				"Scans will begin immediately on online devices. Results will be reported to Microsoft Intune admin center. "+
				"Devices will be protected during scanning and threats will be quarantined automatically.",
				successCount, quickScanCount, fullScanCount),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

func (a *WindowsDefenderScanAction) scanManagedDevice(ctx context.Context, deviceID string, quickScan bool) error {
	scanType := "full"
	if quickScan {
		scanType = "quick"
	}
	tflog.Debug(ctx, fmt.Sprintf("Initiating %s Windows Defender scan on managed device with ID: %s", scanType, deviceID))

	// Construct the request body
	requestBody := constructManagedDeviceRequest(ctx, quickScan)

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderScan().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate %s scan on managed device %s: %w", scanType, deviceID, err)
	}

	return nil
}

func (a *WindowsDefenderScanAction) scanComanagedDevice(ctx context.Context, deviceID string, quickScan bool) error {
	scanType := "full"
	if quickScan {
		scanType = "quick"
	}
	tflog.Debug(ctx, fmt.Sprintf("Initiating %s Windows Defender scan on co-managed device with ID: %s", scanType, deviceID))

	// Construct the request body
	requestBody := constructComanagedDeviceRequest(ctx, quickScan)

	err := a.client.
		DeviceManagement().
		ComanagedDevices().
		ByManagedDeviceId(deviceID).
		WindowsDefenderScan().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("failed to initiate %s scan on co-managed device %s: %w", scanType, deviceID, err)
	}

	return nil
}
