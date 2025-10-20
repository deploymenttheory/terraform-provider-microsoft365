package graphBetaDeleteUserFromSharedAppleDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (a *DeleteUserFromSharedAppleDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data DeleteUserFromSharedAppleDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check that at least one device list is provided
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device-user pair.",
		)
		return
	}

	// Check for duplicate device-user pairs within managed devices
	if len(data.ManagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, deviceUser := range data.ManagedDevices {
			deviceID := deviceUser.DeviceID.ValueString()
			upn := deviceUser.UserPrincipalName.ValueString()
			key := fmt.Sprintf("%s:%s", deviceID, upn)
			if seen[key] {
				duplicates = append(duplicates, fmt.Sprintf("Device: %s, User: %s", deviceID, upn))
			}
			seen[key] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Duplicate Device-User Pairs Found",
				fmt.Sprintf("The following device-user pairs are duplicated in managed_devices: %s. "+
					"Each user will only be deleted once from each device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, "; ")),
			)
		}
	}

	// Check for duplicate device-user pairs within co-managed devices
	if len(data.ComanagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, deviceUser := range data.ComanagedDevices {
			deviceID := deviceUser.DeviceID.ValueString()
			upn := deviceUser.UserPrincipalName.ValueString()
			key := fmt.Sprintf("%s:%s", deviceID, upn)
			if seen[key] {
				duplicates = append(duplicates, fmt.Sprintf("Device: %s, User: %s", deviceID, upn))
			}
			seen[key] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_devices"),
				"Duplicate Device-User Pairs Found",
				fmt.Sprintf("The following device-user pairs are duplicated in comanaged_devices: %s. "+
					"Each user will only be deleted once from each device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, "; ")),
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating delete user from shared Apple device action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentDevices []string
	var nonIPadDevices []string
	var unsupervisedDevices []string
	deviceMap := make(map[string]bool) // Track unique devices to avoid duplicate API calls

	// Validate managed devices
	for _, deviceUser := range data.ManagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()

		// Skip if we've already validated this device
		if deviceMap[deviceID] {
			continue
		}
		deviceMap[deviceID] = true

		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentDevices = append(nonExistentDevices, fmt.Sprintf("%s (managed)", deviceID))
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("managed_devices"),
					"Error Validating Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
				return
			}
		} else {
			// Check that device is iPad (iPadOS)
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				// Only iPad devices support Shared iPad mode
				if deviceType != models.IPAD_DEVICETYPE {
					nonIPadDevices = append(nonIPadDevices,
						fmt.Sprintf("%s (managed, deviceType: %s)", deviceID, deviceType.String()))
				} else {
					// For iPad devices, check if supervised
					if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
						unsupervisedDevices = append(unsupervisedDevices, fmt.Sprintf("%s (managed)", deviceID))
					}
				}
			} else {
				nonIPadDevices = append(nonIPadDevices, fmt.Sprintf("%s (managed, Unknown deviceType)", deviceID))
			}
		}
	}

	// Validate co-managed devices
	for _, deviceUser := range data.ComanagedDevices {
		deviceID := deviceUser.DeviceID.ValueString()

		// Skip if we've already validated this device
		if deviceMap[deviceID] {
			continue
		}
		deviceMap[deviceID] = true

		device, err := a.client.
			DeviceManagement().
			ComanagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentDevices = append(nonExistentDevices, fmt.Sprintf("%s (co-managed)", deviceID))
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("comanaged_devices"),
					"Error Validating Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
				return
			}
		} else {
			// Check that device is iPad (iPadOS)
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				// Only iPad devices support Shared iPad mode
				if deviceType != models.IPAD_DEVICETYPE {
					nonIPadDevices = append(nonIPadDevices,
						fmt.Sprintf("%s (co-managed, deviceType: %s)", deviceID, deviceType.String()))
				} else {
					// For iPad devices, check if supervised
					if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
						unsupervisedDevices = append(unsupervisedDevices, fmt.Sprintf("%s (co-managed)", deviceID))
					}
				}
			} else {
				nonIPadDevices = append(nonIPadDevices, fmt.Sprintf("%s (co-managed, Unknown deviceType)", deviceID))
			}
		}
	}

	if len(nonExistentDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("devices"),
			"Non-Existent Devices",
			fmt.Sprintf("The following device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentDevices, ", ")),
		)
	}

	if len(nonIPadDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("devices"),
			"Non-iPad Devices",
			fmt.Sprintf("The delete user from shared Apple device action only works on iPadOS devices in Shared iPad mode. "+
				"The following devices are not iPadOS devices: %s. "+
				"Please remove non-iPadOS devices from the devices list.",
				strings.Join(nonIPadDevices, ", ")),
		)
	}

	if len(unsupervisedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("devices"),
			"Unsupervised iPadOS Devices",
			fmt.Sprintf("The following iPadOS devices are not supervised: %s. "+
				"Shared iPad mode requires supervised devices enrolled via DEP/ABM. "+
				"Unsupervised devices cannot use Shared iPad mode, and the delete user action will fail for these devices.",
				strings.Join(unsupervisedDevices, ", ")),
		)
	}

	// Critical warning about data deletion
	resp.Diagnostics.AddAttributeWarning(
		path.Root("devices"),
		"Permanent Data Deletion Warning",
		fmt.Sprintf("This action will permanently delete %d user(s) and all their cached data from their respective Shared iPad device(s). "+
			"This operation CANNOT be undone. Deleted users will need to be re-added if they need access to these devices again, "+
			"and they will lose all locally cached data (documents, photos, app data). "+
			"The user accounts in the cloud (Azure AD/Entra ID) will NOT be affected, only their cached presence on these specific devices.",
			len(data.ManagedDevices)+len(data.ComanagedDevices)),
	)

	// Warning about Shared iPad mode requirement
	resp.Diagnostics.AddAttributeWarning(
		path.Root("devices"),
		"Shared iPad Mode Requirement",
		fmt.Sprintf("This action only works on iPads configured in Shared iPad mode. "+
			"Regular (non-shared) iPads will not be affected by this action. "+
			"Ensure the %d device(s) in this action are actually configured in Shared iPad mode. "+
			"The action will fail if devices are not in Shared iPad mode or if the specified users do not exist on the devices.",
			len(deviceMap)),
	)
}
