package graphBetaLogoutSharedAppleDeviceActiveUser

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *LogoutSharedAppleDeviceActiveUserAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data LogoutSharedAppleDeviceActiveUserActionModel

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
				"Each device will only have logout performed once, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating logout shared Apple device action for %d device(s)", len(deviceIDs)))

	var nonExistentDevices []string
	var nonIPadDevices []string
	var unsupervisedDevices []string
	var potentiallyNotSharedDevices []string

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
			// Check that device is iPadOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if os != "ipados" && os != "ios" {
					nonIPadDevices = append(nonIPadDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				} else {
					// For iPadOS/iOS devices, check if supervised
					if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
						unsupervisedDevices = append(unsupervisedDevices, deviceID)
					}

					// Note: We cannot directly verify if a device is in Shared iPad mode from the API
					// The action will simply fail gracefully if the device is not in Shared iPad mode
					// We can check for iOS devices (iPhones don't support Shared mode)
					if os == "ios" {
						potentiallyNotSharedDevices = append(potentiallyNotSharedDevices,
							fmt.Sprintf("%s (iOS - likely iPhone, not Shared iPad)", deviceID))
					}
				}
			} else {
				nonIPadDevices = append(nonIPadDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
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

	if len(nonIPadDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Non-iPad Devices",
			fmt.Sprintf("The logout shared Apple device action only works on iPadOS devices in Shared iPad mode. "+
				"The following devices are not iPadOS devices: %s. "+
				"Please remove non-iPadOS devices from the device_ids list.",
				strings.Join(nonIPadDevices, ", ")),
		)
	}

	if len(potentiallyNotSharedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Potential Non-Shared iPad Devices",
			fmt.Sprintf("The following devices may not support Shared iPad mode: %s. "+
				"Shared iPad mode is only available on iPadOS devices (not iPhones). "+
				"The action may fail for these devices.",
				strings.Join(potentiallyNotSharedDevices, ", ")),
		)
	}

	if len(unsupervisedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Unsupervised iPadOS Devices",
			fmt.Sprintf("The following iPadOS devices are not supervised: %s. "+
				"Shared iPad mode requires supervised devices enrolled via DEP/ABM. "+
				"Unsupervised devices cannot use Shared iPad mode, and the logout action will fail for these devices.",
				strings.Join(unsupervisedDevices, ", ")),
		)
	}

	// General warning about Shared iPad mode requirement
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"Shared iPad Mode Requirement",
		fmt.Sprintf("This action only works on iPads configured in Shared iPad mode. "+
			"Regular (non-shared) iPads will not be affected by this action, even if they meet other requirements. "+
			"Ensure the %d device(s) in this action are actually configured in Shared iPad mode. "+
			"The action will fail gracefully if devices are not in Shared iPad mode or if no user is currently logged in.",
			len(deviceIDs)),
	)
}
