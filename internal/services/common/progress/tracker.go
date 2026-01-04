package progress

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// For creates a new progress tracker for action invocation
// Usage: progress.For(resp).Starting("action").Device(id).Succeeded("message")
func For(resp *action.InvokeResponse) *progressTracker {
	return &progressTracker{
		resp: resp,
	}
}

// progressTracker tracks action progress and sends progress events
type progressTracker struct {
	resp         *action.InvokeResponse
	successCount int
	failureCount int
	totalDevices int
}

// WithTotalDevices sets the total device count for final messages
func (p *progressTracker) WithTotalDevices(total int) *progressTracker {
	p.totalDevices = total
	return p
}

// Starting sends an initial progress message
func (p *progressTracker) Starting(actionName string, details ...string) *progressTracker {
	var message string
	if len(details) > 0 {
		message = fmt.Sprintf("Starting %s for %d device(s) (%s)",
			actionName, p.totalDevices, details[0])
	} else {
		message = fmt.Sprintf("Starting %s for %d device(s)",
			actionName, p.totalDevices)
	}

	p.resp.SendProgress(action.InvokeProgressEvent{
		Message: message,
	})
	return p
}

// Device starts a device status chain
func (p *progressTracker) Device(deviceID string, deviceType ...string) *deviceProgress {
	dt := ""
	if len(deviceType) > 0 {
		dt = deviceType[0]
	}
	return &deviceProgress{
		tracker:    p,
		deviceID:   deviceID,
		deviceType: dt,
	}
}

// Message sends a custom progress message
func (p *progressTracker) Message(message string) *progressTracker {
	p.resp.SendProgress(action.InvokeProgressEvent{
		Message: message,
	})
	return p
}

// CompletedWithIgnoredFailures reports completion with ignored failures
func (p *progressTracker) CompletedWithIgnoredFailures(actionName string) *progressTracker {
	p.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("%s completed: %d succeeded, %d failed (ignored)",
			actionName, p.successCount, p.failureCount),
	})
	return p
}

// CompletedSuccessfully reports successful completion
func (p *progressTracker) CompletedSuccessfully(message string) *progressTracker {
	p.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully %s on all %d device(s)",
			message, p.totalDevices),
	})
	return p
}

// Failed adds a diagnostic error with failure counts
func (p *progressTracker) Failed(title string, action string) {
	p.resp.Diagnostics.AddError(
		title,
		fmt.Sprintf("Failed to %s on %d of %d device(s)",
			action, p.failureCount, p.totalDevices),
	)
}

// SuccessCount returns the number of successful operations
func (p *progressTracker) SuccessCount() int {
	return p.successCount
}

// FailureCount returns the number of failed operations
func (p *progressTracker) FailureCount() int {
	return p.failureCount
}

// HasFailures returns true if any operations failed
func (p *progressTracker) HasFailures() bool {
	return p.failureCount > 0
}

// deviceProgress reports device operation status
type deviceProgress struct {
	tracker    *progressTracker
	deviceID   string
	deviceType string
}

// Succeeded reports device operation succeeded
func (d *deviceProgress) Succeeded(message string) *progressTracker {
	d.tracker.successCount++

	prefix := "Device"
	if d.deviceType != "" {
		prefix = d.deviceType + " device"
	}

	d.tracker.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("%s %s: %s", prefix, d.deviceID, message),
	})
	return d.tracker
}

// Failed reports device operation failed
func (d *deviceProgress) Failed(errorMessage string) *progressTracker {
	d.tracker.failureCount++

	prefix := "Device"
	if d.deviceType != "" {
		prefix = d.deviceType + " device"
	}

	d.tracker.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("%s %s: %s", prefix, d.deviceID, errorMessage),
	})
	return d.tracker
}
