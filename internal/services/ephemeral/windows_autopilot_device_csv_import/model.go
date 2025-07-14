package windowsAutopilotDeviceCSVImport

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsAutopilotDeviceCSVImportModel represents the ephemeral resource model for importing Windows Autopilot devices from a CSV file
type WindowsAutopilotDeviceCSVImportModel struct {
	// Path to the CSV file containing device information
	FilePath types.String `tfsdk:"file_path"`

	// Devices imported from the CSV file
	Devices []DeviceEntry `tfsdk:"devices"`
}

// DeviceEntry represents a single device entry from a CSV file
type DeviceEntry struct {
	// Device Serial Number (required)
	SerialNumber types.String `tfsdk:"serial_number"`

	// Windows Product ID (optional for admins, required for partners)
	WindowsProductID types.String `tfsdk:"windows_product_id"`

	// Hardware Hash (required)
	HardwareHash types.String `tfsdk:"hardware_hash"`

	// Group Tag (optional)
	GroupTag types.String `tfsdk:"group_tag"`

	// Assigned User (optional)
	AssignedUser types.String `tfsdk:"assigned_user"`
}
