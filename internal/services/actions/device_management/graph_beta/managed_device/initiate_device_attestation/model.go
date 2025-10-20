// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatedeviceattestation?view=graph-rest-beta
package graphBetaInitiateDeviceAttestationManagedDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InitiateDeviceAttestationManagedDeviceActionModel struct {
	ManagedDeviceIDs   types.List     `tfsdk:"managed_device_ids"`
	ComanagedDeviceIDs types.List     `tfsdk:"comanaged_device_ids"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

