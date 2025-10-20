package graphBetaBypassActivationLockManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (a *BypassActivationLockManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data BypassActivationLockManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deviceIDs := make([]string, 0, len(data.DeviceIDs.Elements()))
	for _, elem := range data.DeviceIDs.Elements() {
		deviceIDs = append(deviceIDs, elem.String())
	}

	seen := make(map[string]bool)
	var duplicates []string
	for _, id := range deviceIDs {
		if seen[id] {
			duplicates = append(duplicates, id)
		}
		seen[id] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs Found",
			fmt.Sprintf("The following device IDs are duplicated in the configuration: %s. "+
				"Activation Lock bypass will only be issued once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating bypass activation lock action for %d device(s)", len(deviceIDs)))

	// Only perform online validation if enabled and client is configured
	validateExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateExists = data.ValidateDeviceExists.ValueBool()
	}

	if !validateExists {
		tflog.Debug(ctx, "Device existence validation disabled, skipping online validation")
		return
	}

	if a.client == nil {
		tflog.Debug(ctx, "Client not configured, skipping device existence validation")
		return
	}

	var nonExistentDevices []string
	var unsupportedOSDevices []string
	var unsupervisedIOSDevices []string
	var activationLockDisabledDevices []string
	var alreadyBypassedDevices []string

	for _, deviceID := range deviceIDs {
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentDevices = append(nonExistentDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("device_ids"),
					"Error Validating Device Existence",
					fmt.Sprintf("Failed to check existence of device %s: %s", deviceID, err.Error()),
				)
				return
			}
		} else {

			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				// Activation Lock bypass is supported on iPhone, iPad, and Mac
				supportsActivationLockBypass := deviceType == models.IPHONE_DEVICETYPE ||
					deviceType == models.IPAD_DEVICETYPE ||
					deviceType == models.MAC_DEVICETYPE

				if !supportsActivationLockBypass {
					unsupportedOSDevices = append(unsupportedOSDevices,
						fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
					continue
				}

				// iOS/iPadOS must be supervised
				if (deviceType == models.IPHONE_DEVICETYPE || deviceType == models.IPAD_DEVICETYPE) &&
					(device.GetIsSupervised() == nil || !*device.GetIsSupervised()) {
					unsupervisedIOSDevices = append(unsupervisedIOSDevices,
						fmt.Sprintf("%s (deviceType: %s - not supervised)", deviceID, deviceType.String()))
					continue
				}
			} else {
				unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
				continue
			}

			// Check if Activation Lock is actually enabled
			// The activationLockBypassCode field presence indicates bypass capability
			// If the field is already populated, device may already have bypass code
			if device.GetActivationLockBypassCode() != nil && *device.GetActivationLockBypassCode() != "" {
				alreadyBypassedDevices = append(alreadyBypassedDevices, deviceID)
			}

			// Check if device has Activation Lock enabled
			// Note: Not all API versions expose this directly, so we check via proxy indicators
			// Supervised iOS/iPadOS and DEP macOS with Find My enabled should have Activation Lock
			deviceType := *device.GetDeviceType()
			if (deviceType == models.IPHONE_DEVICETYPE || deviceType == models.IPAD_DEVICETYPE) &&
				device.GetIsSupervised() != nil && *device.GetIsSupervised() {
				// For supervised devices, we assume Activation Lock may be enabled
				// The actual check happens when the bypass command is issued
			} else if deviceType == models.MAC_DEVICETYPE {
				// For macOS, if not clearly DEP enrolled, warn user
				activationLockDisabledDevices = append(activationLockDisabledDevices,
					fmt.Sprintf("%s (may not have Activation Lock enabled)", deviceID))
			}
		}
	}

	if len(nonExistentDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Non-Existent Devices",
			fmt.Sprintf("The following device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentDevices, ", ")),
		)
	}

	if len(unsupportedOSDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Unsupported Devices for Activation Lock Bypass",
			fmt.Sprintf("Activation Lock bypass is only supported on iOS, iPadOS, and macOS devices. The following devices are not supported: %s. "+
				"Please remove non-Apple devices from the configuration.",
				strings.Join(unsupportedOSDevices, ", ")),
		)
	}

	if len(unsupervisedIOSDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Unsupervised iOS/iPadOS Devices",
			fmt.Sprintf("iOS and iPadOS devices must be supervised to support Activation Lock bypass. The following devices are not supervised: %s. "+
				"Supervise these devices via DEP/ABM enrollment or Apple Configurator before attempting bypass.",
				strings.Join(unsupervisedIOSDevices, ", ")),
		)
	}

	if len(alreadyBypassedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Devices Already Have Bypass Codes",
			fmt.Sprintf("The following devices already have Activation Lock bypass codes: %s. "+
				"Issuing a new bypass command will generate a new code, but the existing code may still be valid. "+
				"Check device properties in Intune portal to retrieve existing bypass codes before generating new ones.",
				strings.Join(alreadyBypassedDevices, ", ")),
		)
	}

	if len(activationLockDisabledDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Devices May Not Have Activation Lock Enabled",
			fmt.Sprintf("The following devices may not have Activation Lock enabled: %s. "+
				"Activation Lock is only active when Find My iPhone/iPad/Mac is enabled by the user. "+
				"Verify Activation Lock status before attempting bypass, as the command may fail if Activation Lock is not enabled.",
				strings.Join(activationLockDisabledDevices, ", ")),
		)
	}
}
