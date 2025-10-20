package graphBetaRecoverManagedDevicePasscode

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RecoverManagedDevicePasscodeAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RecoverManagedDevicePasscodeActionModel

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
				"Passcode recovery will only be attempted once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating passcode recovery action for %d device(s)", len(deviceIDs)))

	var nonExistentDevices []string
	var unsupportedOSDevices []string
	var unsupervisedDevices []string
	var escrowWarningDevices []string

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
			// Check OS compatibility - Passcode recovery primarily works on iOS/iPadOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if os != "ios" && os != "ipados" {
					unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				} else {
					// For iOS/iPadOS, check if supervised (better passcode escrow support)
					if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
						unsupervisedDevices = append(unsupervisedDevices, deviceID)
					}
				}
			} else {
				unsupportedOSDevices = append(unsupportedOSDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			}

			// Note: We cannot directly check if a passcode is escrowed from device properties
			// This would require checking the actual passcode field which may not be exposed
			// We'll add a general warning about escrow requirements
			escrowWarningDevices = append(escrowWarningDevices, deviceID)
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
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Devices with Limited Passcode Recovery Support",
			fmt.Sprintf("The following devices may have limited or no support for passcode recovery: %s. "+
				"Passcode recovery is primarily supported on supervised iOS and iPadOS devices. "+
				"The action may fail for these devices if passcode escrow is not supported or configured.",
				strings.Join(unsupportedOSDevices, ", ")),
		)
	}

	if len(unsupervisedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Unsupervised iOS/iPadOS Devices",
			fmt.Sprintf("The following iOS/iPadOS devices are not supervised: %s. "+
				"Passcode recovery works best with supervised devices enrolled via DEP/ABM. "+
				"Unsupervised devices may not have passcode escrow enabled.",
				strings.Join(unsupervisedDevices, ", ")),
		)
	}

	if len(escrowWarningDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Passcode Escrow Requirement",
			fmt.Sprintf("Passcode recovery requires that passcodes were escrowed during device enrollment. "+
				"If passcodes were not escrowed for the %d device(s) in this action, recovery will fail. "+
				"Check device enrollment profiles to ensure passcode escrow is enabled. "+
				"If recovery fails, consider using the reset passcode action instead, which generates a new temporary passcode.",
				len(escrowWarningDevices)),
		)
	}
}
