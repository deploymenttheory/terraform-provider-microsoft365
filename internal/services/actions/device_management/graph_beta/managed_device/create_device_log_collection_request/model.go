// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta
package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CreateDeviceLogCollectionRequestManagedDeviceActionModel struct {
	ManagedDevices        []ManagedDeviceLogCollection   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceLogCollection `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                     `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                     `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value                 `tfsdk:"timeouts"`
}

type ManagedDeviceLogCollection struct {
	DeviceID     types.String `tfsdk:"device_id"`
	TemplateType types.String `tfsdk:"template_type"`
}

type ComanagedDeviceLogCollection struct {
	DeviceID     types.String `tfsdk:"device_id"`
	TemplateType types.String `tfsdk:"template_type"`
}
