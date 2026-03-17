package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchDeviceRegistrationResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	UpdateCategory types.String   `tfsdk:"update_category"`
	DeviceIds      types.Set      `tfsdk:"device_ids"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}
