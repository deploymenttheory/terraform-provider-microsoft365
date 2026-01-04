// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-updatewindowsdeviceaccount?view=graph-rest-beta
package graphBetaUpdateWindowsDeviceAccount

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpdateWindowsDeviceAccountActionModel struct {
	ManagedDevices        []ManagedDeviceAccount   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceAccount `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool               `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool               `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value           `tfsdk:"timeouts"`
}

type ManagedDeviceAccount struct {
	DeviceID                         types.String `tfsdk:"device_id"`
	DeviceAccountEmail               types.String `tfsdk:"device_account_email"`
	Password                         types.String `tfsdk:"password"`
	PasswordRotationEnabled          types.Bool   `tfsdk:"password_rotation_enabled"`
	CalendarSyncEnabled              types.Bool   `tfsdk:"calendar_sync_enabled"`
	ExchangeServer                   types.String `tfsdk:"exchange_server"`
	SessionInitiationProtocolAddress types.String `tfsdk:"session_initiation_protocol_address"`
}

type ComanagedDeviceAccount struct {
	DeviceID                         types.String `tfsdk:"device_id"`
	DeviceAccountEmail               types.String `tfsdk:"device_account_email"`
	Password                         types.String `tfsdk:"password"`
	PasswordRotationEnabled          types.Bool   `tfsdk:"password_rotation_enabled"`
	CalendarSyncEnabled              types.Bool   `tfsdk:"calendar_sync_enabled"`
	ExchangeServer                   types.String `tfsdk:"exchange_server"`
	SessionInitiationProtocolAddress types.String `tfsdk:"session_initiation_protocol_address"`
}
