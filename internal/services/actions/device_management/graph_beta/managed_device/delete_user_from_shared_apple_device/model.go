// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta
package graphBetaDeleteUserFromSharedAppleDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeleteUserFromSharedAppleDeviceActionModel struct {
	ManagedDevices        []ManagedDeviceUserPair   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceUserPair `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value            `tfsdk:"timeouts"`
}

type ManagedDeviceUserPair struct {
	DeviceID          types.String `tfsdk:"device_id"`
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

type ComanagedDeviceUserPair struct {
	DeviceID          types.String `tfsdk:"device_id"`
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}
