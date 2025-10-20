package graphBetaRevokeAppleVppLicenses

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RevokeAppleVppLicensesAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RevokeAppleVppLicensesActionModel

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

	// Validate that at least one device list is provided
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
					"Licenses will only be revoked once per device, but you should remove duplicates from your configuration.",
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
					"Licenses will only be revoked once per device, but you should remove duplicates from your configuration.",
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
						"A device should only be in one list. The license revocation will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating Apple VPP license revocation action for %d managed and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	var nonExistentManagedDevices []string
	var nonAppleManagedDevices []string
	var nonExistentComanagedDevices []string
	var nonAppleComanagedDevices []string

	// Validate managed devices exist and are iOS/iPadOS
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
			// Check that device is iOS or iPadOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if os != "ios" && os != "ipados" {
					nonAppleManagedDevices = append(nonAppleManagedDevices,
						fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	// Validate co-managed devices exist and are iOS/iPadOS
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
			// Check that device is iOS or iPadOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				if os != "ios" && os != "ipados" {
					nonAppleComanagedDevices = append(nonAppleComanagedDevices,
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

	if len(nonAppleManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_device_ids"),
			"Non-Apple Devices",
			fmt.Sprintf("The Apple VPP license revocation action only works on iOS and iPadOS devices. "+
				"The following managed devices are not iOS/iPadOS devices: %s. "+
				"Please remove non-Apple devices from the managed_device_ids list.",
				strings.Join(nonAppleManagedDevices, ", ")),
		)
	}

	if len(nonAppleComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Non-Apple Co-Managed Devices",
			fmt.Sprintf("The Apple VPP license revocation action only works on iOS and iPadOS devices. "+
				"The following co-managed devices are not iOS/iPadOS devices: %s. "+
				"Please remove non-Apple devices from the comanaged_device_ids list.",
				strings.Join(nonAppleComanagedDevices, ", ")),
		)
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)

	// Critical warning about license revocation
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"Apple VPP License Revocation Warning",
		fmt.Sprintf("This action will revoke ALL Apple Volume Purchase Program (VPP) licenses from %d device(s). "+
			"This action: "+
			"(1) Returns all VPP app licenses to the available pool "+
			"(2) May remove VPP apps from the devices "+
			"(3) Makes licenses immediately available for reassignment "+
			"(4) Cannot be easily undone (licenses must be manually reassigned) "+
			"(5) Affects all VPP-purchased apps on the devices. "+
			"Only proceed if you intend to reclaim licenses for reallocation or device retirement.",
			totalDevices),
	)

	// Informational message about VPP
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"VPP License Management Information",
		fmt.Sprintf("This action revokes Apple Volume Purchase Program (VPP) licenses from %d iOS/iPadOS device(s). "+
			"VPP licenses are used for apps purchased through Apple Business Manager. "+
			"After revocation, licenses are returned to your organization's available pool and can be viewed in Apple Business Manager. "+
			"Devices must be online to receive app removal commands. "+
			"User data for apps is typically preserved on the device unless the apps are removed by policy.",
			totalDevices),
	)
}
