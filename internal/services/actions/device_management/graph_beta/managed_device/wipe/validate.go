package graphBetaWipeManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *WipeManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data WipeManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.DeviceIDs.IsNull() || data.DeviceIDs.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Missing Required Configuration",
			"device_ids must be specified and contain at least one device ID.",
		)
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(deviceIDs) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Invalid Configuration",
			"device_ids must contain at least one device ID.",
		)
		return
	}

	deviceIDMap := make(map[string]bool)
	var duplicates []string

	for _, deviceID := range deviceIDs {
		if deviceIDMap[deviceID] {
			duplicates = append(duplicates, deviceID)
		}
		deviceIDMap[deviceID] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs",
			fmt.Sprintf("The following device IDs are duplicated in the list: %v. "+
				"Each device will only be wiped once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	// Validate conflicting options
	if !data.KeepUserData.IsNull() && !data.KeepUserData.IsUnknown() && data.KeepUserData.ValueBool() {
		if !data.ObliterationBehavior.IsNull() && !data.ObliterationBehavior.IsUnknown() {
			obliterationBehaviorValue := data.ObliterationBehavior.ValueString()
			if obliterationBehaviorValue != "doNotObliterate" {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("obliteration_behavior"),
					"Conflicting Configuration",
					"When keep_user_data is true, obliteration_behavior should typically be set to 'doNotObliterate'. "+
						"Using obliteration with keep_user_data=true may result in unexpected behavior.",
				)
			}
		}
	}

	// Validate macOS unlock code is only used when needed
	if !data.MacOsUnlockCode.IsNull() && !data.MacOsUnlockCode.IsUnknown() {
		macOsUnlockCode := data.MacOsUnlockCode.ValueString()
		if len(macOsUnlockCode) != 6 {
			resp.Diagnostics.AddAttributeError(
				path.Root("macos_unlock_code"),
				"Invalid macOS Unlock Code",
				"The macOS unlock code must be exactly 6 digits.",
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating wipe action for %d device(s)", len(deviceIDs)))

	var unsupportedDevices []string

	for _, deviceID := range deviceIDs {
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("device_ids"),
				"Unable to Validate Device",
				fmt.Sprintf("Could not fetch managed device with ID %s to validate. "+
					"Ensure the device exists and you have permission to manage it. "+
					"Error: %s", deviceID, err.Error()),
			)
			continue
		}

		if device == nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("device_ids"),
				"Device Not Found",
				fmt.Sprintf("Managed device with ID %s was not found. "+
					"Ensure the device ID is correct and the device is enrolled in Intune.", deviceID),
			)
			continue
		}

		// Check platform compatibility - wipe is supported on Windows, iOS, iPadOS, macOS, and Android (NOT ChromeOS)
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			supportedOS := map[string]bool{
				"windows": true,
				"ios":     true,
				"ipados":  true,
				"macos":   true,
				"android": true,
			}
			if !supportedOS[os] {
				unsupportedDevices = append(unsupportedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				continue
			}
		} else {
			unsupportedDevices = append(unsupportedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		// Check for activation lock on iOS/macOS devices
		operatingSystem := device.GetOperatingSystem()
		if operatingSystem != nil {
			os := *operatingSystem
			if (os == "iOS" || os == "iPadOS" || os == "macOS") &&
				(data.MacOsUnlockCode.IsNull() || data.MacOsUnlockCode.IsUnknown()) {
				activationLockBypassCode := device.GetActivationLockBypassCode()
				if activationLockBypassCode != nil && *activationLockBypassCode != "" {
					resp.Diagnostics.AddAttributeWarning(
						path.Root("macos_unlock_code"),
						"Activation Lock May Be Enabled",
						fmt.Sprintf("Device %s (%s) may have Activation Lock enabled. "+
							"If wiping fails, you may need to provide the macos_unlock_code parameter. "+
							"The bypass code is available in the device details.", deviceID, os),
					)
				}
			}
		}

		// Log device details for debugging
		deviceName := "unknown"
		if device.GetDeviceName() != nil {
			deviceName = *device.GetDeviceName()
		}
		tflog.Debug(ctx, fmt.Sprintf("Validated device %s (Name: %s) for wipe", deviceID, deviceName))
	}

	if len(unsupportedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Unsupported Devices for Wipe",
			fmt.Sprintf("Wipe is supported on Windows, iOS, iPadOS, macOS, and Android devices only. The following devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedDevices, ", ")),
		)
	}

	// Final warning about data loss
	if !data.KeepUserData.IsNull() && !data.KeepUserData.IsUnknown() && !data.KeepUserData.ValueBool() {
		resp.Diagnostics.AddWarning(
			"Data Loss Warning",
			fmt.Sprintf("You are about to wipe %d device(s) with keep_user_data=false. "+
				"This will permanently delete ALL data on these devices. "+
				"This action cannot be undone. Please ensure you have reviewed the device list carefully.", len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation complete for %d device(s)", len(deviceIDs)))
}
