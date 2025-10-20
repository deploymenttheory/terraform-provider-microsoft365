package graphBetaRotateBitLockerKeys

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RotateBitLockerKeysAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RotateBitLockerKeysActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	// Get managed device IDs
	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
	}

	// Get co-managed device IDs
	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if len(managedDeviceIDs) == 0 && len(comanagedDeviceIDs) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_device_ids' or 'comanaged_device_ids' must be provided with at least one device ID.",
		)
		return
	}

	if len(managedDeviceIDs) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, id := range managedDeviceIDs {
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_device_ids"),
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"BitLocker keys will only be rotated once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	if len(comanagedDeviceIDs) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, id := range comanagedDeviceIDs {
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_device_ids"),
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"BitLocker keys will only be rotated once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	for _, managedID := range managedDeviceIDs {
		for _, comanagedID := range comanagedDeviceIDs {
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_device_ids"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_device_ids and comanaged_device_ids. "+
						"A device should only be in one list. The key rotation will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating BitLocker key rotation action for %d managed and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	var nonExistentManagedDevices []string
	var nonWindowsManagedDevices []string
	var nonExistentComanagedDevices []string
	var nonWindowsComanagedDevices []string

	for _, deviceID := range managedDeviceIDs {
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentManagedDevices = append(nonExistentManagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("managed_device_ids"),
					"Error Validating Managed Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check that device is Windows
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if !strings.Contains(os, "windows") {
					nonWindowsManagedDevices = append(nonWindowsManagedDevices,
						fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	for _, deviceID := range comanagedDeviceIDs {
		device, err := a.client.
			DeviceManagement().
			ComanagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentComanagedDevices = append(nonExistentComanagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("comanaged_device_ids"),
					"Error Validating Co-Managed Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check that device is Windows
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if !strings.Contains(os, "windows") {
					nonWindowsComanagedDevices = append(nonWindowsComanagedDevices,
						fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully", deviceID))
		}
	}

	if len(nonExistentManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_device_ids"),
			"Non-Existent Managed Devices",
			fmt.Sprintf("The following managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentManagedDevices, ", ")),
		)
	}

	if len(nonExistentComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Non-Existent Co-Managed Devices",
			fmt.Sprintf("The following co-managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing co-managed devices.",
				strings.Join(nonExistentComanagedDevices, ", ")),
		)
	}

	if len(nonWindowsManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_device_ids"),
			"Non-Windows Devices",
			fmt.Sprintf("The BitLocker key rotation action only works on Windows devices. "+
				"The following managed devices are not Windows devices: %s. "+
				"Please remove non-Windows devices from the managed_device_ids list.",
				strings.Join(nonWindowsManagedDevices, ", ")),
		)
	}

	if len(nonWindowsComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Non-Windows Co-Managed Devices",
			fmt.Sprintf("The BitLocker key rotation action only works on Windows devices. "+
				"The following co-managed devices are not Windows devices: %s. "+
				"Please remove non-Windows devices from the comanaged_device_ids list.",
				strings.Join(nonWindowsComanagedDevices, ", ")),
		)
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)

	// Critical warning about key rotation
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"BitLocker Key Rotation Warning",
		fmt.Sprintf("This action will rotate BitLocker recovery keys on %d Windows device(s). "+
			"This action: "+
			"(1) Generates new BitLocker recovery passwords "+
			"(2) Escrows new keys to Intune/Azure AD "+
			"(3) INVALIDATES all previous recovery keys "+
			"(4) Makes previous keys unusable for device recovery "+
			"(5) Requires devices to be online and connected. "+
			"After rotation, only the new recovery keys can be used to unlock devices if BitLocker recovery is triggered. "+
			"Ensure you have documented processes to retrieve the new keys from Intune/Azure AD when needed.",
			totalDevices),
	)

	// Informational message about BitLocker
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"BitLocker Key Rotation Information",
		fmt.Sprintf("This action rotates BitLocker recovery keys on %d Windows device(s). "+
			"BitLocker recovery keys are used to unlock drives when normal unlock methods fail (forgotten password, TPM issues, etc.). "+
			"The key rotation process: "+
			"(1) Generates new recovery passwords on each device "+
			"(2) Escrows the new keys to Intune and Azure AD "+
			"(3) Invalidates previous recovery keys "+
			"(4) Does NOT re-encrypt data (fast operation) "+
			"(5) Does NOT require device restart "+
			"(6) Does NOT require user interaction. "+
			"New recovery keys can be retrieved from the Intune admin center (Devices > All devices > [device] > Recovery keys) "+
			"or Azure AD portal (Devices > All devices > [device] > BitLocker keys). "+
			"Devices must have BitLocker enabled and properly configured for this action to succeed.",
			totalDevices),
	)
}
