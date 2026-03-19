package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchDeviceRegistrationResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	UpdateCategory       types.String   `tfsdk:"update_category"`
	EntraDeviceObjectIds types.Set      `tfsdk:"entra_device_object_ids"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
