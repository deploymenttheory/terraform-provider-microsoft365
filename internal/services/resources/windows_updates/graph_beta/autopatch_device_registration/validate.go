package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// validateRequest performs validation before enrolling or unenrolling devices.
// It orchestrates validation by calling atomic validation functions.
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) validateRequest(
	ctx context.Context,
	data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel,
	diagnostics *diag.Diagnostics,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Starting device validation for %s resource", ResourceName))

	if data.EntraDeviceObjectIds.IsNull() || data.EntraDeviceObjectIds.IsUnknown() {
		tflog.Debug(ctx, "No device IDs to validate")
		return nil
	}

	elements := data.EntraDeviceObjectIds.Elements()
	if len(elements) == 0 {
		tflog.Debug(ctx, "Empty device IDs set")
		return nil
	}

	deviceIDs := make([]string, 0, len(elements))
	for _, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			continue
		}
		if deviceID := strVal.ValueString(); deviceID != "" {
			deviceIDs = append(deviceIDs, deviceID)
		}
	}

	if err := r.validateDevicesEligibleForEnrollment(ctx, deviceIDs, diagnostics); err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Device validation completed successfully for %d devices", len(deviceIDs)))
	return nil
}

// validateDevicesEligibleForEnrollment validates that all devices are eligible for Windows Updates enrollment.
// This is the atomic validation function that checks the updatable assets collection.
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) validateDevicesEligibleForEnrollment(
	ctx context.Context,
	deviceIDs []string,
	diagnostics *diag.Diagnostics,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating %d devices for Windows Updates enrollment eligibility", len(deviceIDs)))

	filter := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
	result, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
			QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		})

	if err != nil {
		tflog.Error(ctx, "Failed to query updatable assets for validation", map[string]any{
			"error": err.Error(),
		})
		diagnostics.AddError(
			"Failed to Validate Devices",
			fmt.Sprintf("Could not query updatable assets to validate devices: %s", err.Error()),
		)
		return fmt.Errorf("failed to query updatable assets: %w", err)
	}

	deviceMap := make(map[string]graphmodelswindowsupdates.AzureADDeviceable)
	devices := result.GetValue()
	for _, asset := range devices {
		if asset == nil {
			continue
		}

		azureDevice, ok := asset.(graphmodelswindowsupdates.AzureADDeviceable)
		if !ok {
			continue
		}

		deviceID := azureDevice.GetId()
		if deviceID != nil {
			deviceMap[*deviceID] = azureDevice
		}
	}

	var devicesWithErrors []string
	var devicesNotFound []string

	for _, deviceID := range deviceIDs {
		device, found := deviceMap[deviceID]
		if !found {
			devicesNotFound = append(devicesNotFound, deviceID)
			tflog.Warn(ctx, "Device not found in updatable assets", map[string]any{
				"deviceId": deviceID,
			})
			continue
		}

		deviceErrors := device.GetErrors()
		if len(deviceErrors) > 0 {
			var errorReasons []string
			for _, errObj := range deviceErrors {
				if regErr, ok := errObj.(graphmodelswindowsupdates.AzureADDeviceRegistrationErrorable); ok {
					if reason := regErr.GetReason(); reason != nil {
						errorReasons = append(errorReasons, reason.String())
					}
				}
			}
			devicesWithErrors = append(devicesWithErrors, fmt.Sprintf("%s (%s)", deviceID, strings.Join(errorReasons, ", ")))
			tflog.Error(ctx, "Device has registration errors", map[string]any{
				"deviceId": deviceID,
				"errors":   errorReasons,
			})
		} else {
			tflog.Debug(ctx, "Device is eligible for enrollment", map[string]any{
				"deviceId": deviceID,
			})
		}
	}

	if len(devicesNotFound) > 0 {
		diagnostics.AddError(
			"Devices Not Found in Updatable Assets",
			fmt.Sprintf("The following device IDs are not registered as updatable assets: %s. "+
				"Only devices that are registered with Windows Updates can be enrolled. "+
				"Ensure the devices are Azure AD joined and meet Windows Updates requirements.",
				strings.Join(devicesNotFound, ", ")),
		)
		return fmt.Errorf("%w: %d devices not found in updatable assets", sentinels.ErrEntraDeviceValidationFailed, len(devicesNotFound))
	}

	if len(devicesWithErrors) > 0 {
		diagnostics.AddError(
			"Devices Not Eligible for Enrollment",
			fmt.Sprintf("The following devices have registration errors and cannot be enrolled: %s. "+
				"These devices may be stale, deleted from Entra ID, or have other issues. "+
				"Check the Windows Updates admin center for device status.",
				strings.Join(devicesWithErrors, "; ")),
		)
		return fmt.Errorf("%w: %d devices have registration errors", sentinels.ErrEntraDeviceValidationFailed, len(devicesWithErrors))
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d devices validated successfully", len(deviceIDs)))
	return nil
}
