package windowsAutopilotDeviceCSVImport

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

const (
	EphemeralResourceName = "microsoft365_windows_autopilot_device_csv_import"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ ephemeral.EphemeralResource              = &WindowsAutopilotDeviceCSVImportEphemeralResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &WindowsAutopilotDeviceCSVImportEphemeralResource{}
)

// NewWindowsAutopilotDeviceCSVImportEphemeralResource is a helper function to simplify provider implementation
func NewWindowsAutopilotDeviceCSVImportEphemeralResource() ephemeral.EphemeralResource {
	return &WindowsAutopilotDeviceCSVImportEphemeralResource{}
}

// WindowsAutopilotDeviceCSVImportEphemeralResource is the ephemeral resource implementation
type WindowsAutopilotDeviceCSVImportEphemeralResource struct {
	// This resource doesn't need a Graph client as it only reads from CSV files
}

// Metadata returns the resource type name
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = EphemeralResourceName
}

// Schema defines the schema for the ephemeral resource
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Imports Windows Autopilot devices from a CSV file as an ephemeral resource. This ephemeral resource is used to parse and validate device information from CSV files without persisting in state.",
		Attributes: map[string]schema.Attribute{
			"file_path": schema.StringAttribute{
				MarkdownDescription: "Path to the CSV file containing Windows Autopilot device information. The CSV file must be in ANSI format with no quotation marks and include the required headers: `Device Serial Number`, `Windows Product ID`, `Hardware Hash`, `Group Tag`, and `Assigned User`.",
				Required:            true,
			},
			"devices": schema.ListNestedAttribute{
				MarkdownDescription: "List of devices imported from the CSV file.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"serial_number": schema.StringAttribute{
							MarkdownDescription: "Device Serial Number.",
							Computed:            true,
						},
						"windows_product_id": schema.StringAttribute{
							MarkdownDescription: "Windows Product ID.",
							Computed:            true,
						},
						"hardware_hash": schema.StringAttribute{
							MarkdownDescription: "Hardware Hash.",
							Computed:            true,
						},
						"group_tag": schema.StringAttribute{
							MarkdownDescription: "Group Tag.",
							Computed:            true,
						},
						"assigned_user": schema.StringAttribute{
							MarkdownDescription: "Assigned User.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
