package graphBetaWipeManagedDevice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/progress"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validation"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func (a *WipeManagedDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data WipeManagedDeviceActionModel

	tflog.Debug(ctx, "Starting device wipe", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	totalDevices := len(deviceIDs)
	tflog.Debug(ctx, "Processing devices for wipe", map[string]any{
		"total_devices": totalDevices,
	})

	ignorePartialFailures := false
	if !data.IgnorePartialFailures.IsNull() && !data.IgnorePartialFailures.IsUnknown() {
		ignorePartialFailures = data.IgnorePartialFailures.ValueBool()
	}

	validateDeviceExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateDeviceExists = data.ValidateDeviceExists.ValueBool()
	}

	macOsUnlockCodeProvided := !data.MacOsUnlockCode.IsNull() && !data.MacOsUnlockCode.IsUnknown()

	if validateDeviceExists {
		tflog.Debug(ctx, "Performing device validation via API")

		validationResult, err := validateRequest(ctx, a.client, deviceIDs, macOsUnlockCodeProvided)
		if err != nil {
			tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError("Device Validation Failed", fmt.Sprintf("Failed to validate devices: %s", err.Error()))
			return
		}

		results := validation.NewResults().
			Error(validationResult.NonExistentDevices, "device", "do not exist or are not managed by Intune").
			Error(validationResult.UnsupportedDevices, "device", "are not supported for wipe. Wipe is supported on Windows, iOS, iPadOS, macOS, and Android devices only")

		// Add activation lock warnings
		if len(validationResult.ActivationLockWarningIDs) > 0 {
			for i, deviceID := range validationResult.ActivationLockWarningIDs {
				os := validationResult.ActivationLockWarningOSList[i]
				results = results.Warning([]string{deviceID}, "device", fmt.Sprintf("(%s) may have Activation Lock enabled. If wiping fails, you may need to provide the macos_unlock_code parameter. The bypass code is available in the device details", os))
			}
		}

		if results.Report(resp) {
			return
		}

		tflog.Debug(ctx, "Device validation completed successfully")
	} else {
		tflog.Debug(ctx, "Device validation disabled, skipping API checks")
	}

	// Construct the request body once (same parameters for all devices)
	requestBody, err := constructRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Constructing Request",
			fmt.Sprintf("Could not construct request for wipe managed device: %s", err.Error()),
		)
		return
	}

	progressTracker := progress.For(resp).WithTotalDevices(totalDevices).
		Starting("device wipe operations", fmt.Sprintf("%d devices", totalDevices))

	for _, deviceID := range deviceIDs {
		err := a.wipeDevice(ctx, deviceID, requestBody)
		if err != nil {
			progressTracker.Device(deviceID, "").Failed(fmt.Sprintf("wipe failed: %s", err.Error()))
			tflog.Error(ctx, "Failed to wipe device", map[string]any{"device_id": deviceID, "error": err.Error()})
		} else {
			progressTracker.Device(deviceID, "").Succeeded("wipe initiated successfully")
			tflog.Info(ctx, "Successfully initiated wipe on device", map[string]any{"device_id": deviceID})
		}
	}

	if progressTracker.HasFailures() {
		if ignorePartialFailures {
			progressTracker.CompletedWithIgnoredFailures("device wipe operations")
			tflog.Warn(ctx, "Device wipe completed with ignored failures", map[string]any{
				"success_count": progressTracker.SuccessCount(),
				"failed_count":  progressTracker.FailureCount(),
			})
		} else {
			progressTracker.Failed("Device Wipe Failed", "wipe devices")
			return
		}
	} else {
		// Build completion message based on configuration
		var dataHandling []string
		if !data.KeepUserData.IsNull() && !data.KeepUserData.IsUnknown() && data.KeepUserData.ValueBool() {
			dataHandling = append(dataHandling, "User data will be preserved")
		} else {
			dataHandling = append(dataHandling, "All data (company and personal) will be permanently deleted")
		}

		if !data.KeepEnrollmentData.IsNull() && !data.KeepEnrollmentData.IsUnknown() && data.KeepEnrollmentData.ValueBool() {
			dataHandling = append(dataHandling, "enrollment state maintained for automatic re-enrollment")
		}

		if !data.UseProtectedWipe.IsNull() && !data.UseProtectedWipe.IsUnknown() && data.UseProtectedWipe.ValueBool() {
			dataHandling = append(dataHandling, "protected wipe enabled (UEFI licenses preserved on Windows)")
		}

		completionMsg := "initiated wipe on all devices. " + strings.Join(dataHandling, ", ") + ". This action cannot be undone"
		progressTracker.CompletedSuccessfully(completionMsg)
	}

	tflog.Info(ctx, "Device wipe completed", map[string]any{
		"success_count":            progressTracker.SuccessCount(),
		"failed_count":             progressTracker.FailureCount(),
		"total_devices":            totalDevices,
		"partial_failures_ignored": ignorePartialFailures && progressTracker.HasFailures(),
	})
}

func (a *WipeManagedDeviceAction) wipeDevice(ctx context.Context, deviceID string, requestBody *devicemanagement.ManagedDevicesItemWipePostRequestBody) error {
	tflog.Debug(ctx, "Wiping device", map[string]any{"device_id": deviceID})

	err := a.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceID).
		Wipe().
		Post(ctx, requestBody, nil)

	if err != nil {
		return fmt.Errorf("%s", errors.HandleKiotaGraphErrorForAction(ctx, err))
	}

	return nil
}
