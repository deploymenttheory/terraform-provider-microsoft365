package graphBetaDeleteUserFromSharedAppleDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_devices":   len(data.ManagedDevices),
		"comanaged_devices": len(data.ComanagedDevices),
	})
}
