package windowsAutopilotDeviceCSVImport

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ ephemeral.EphemeralResource = &WindowsAutopilotDeviceCSVImportEphemeralResource{}
)

// NewWindowsAutopilotDeviceCSVImportEphemeralResource is a helper function to simplify provider implementation
func NewWindowsAutopilotDeviceCSVImportEphemeralResource() ephemeral.EphemeralResource {
	return &WindowsAutopilotDeviceCSVImportEphemeralResource{}
}

// WindowsAutopilotDeviceCSVImportEphemeralResource is the ephemeral resource implementation
type WindowsAutopilotDeviceCSVImportEphemeralResource struct {
	// Add any fields needed for the resource implementation
}

// Metadata returns the resource type name
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_autopilot_device_csv_import"
}

// Schema defines the schema for the ephemeral resource
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Imports Windows Autopilot devices from a CSV file. This is an ephemeral resource that does not persist in state.",
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

// Open is called when the ephemeral resource is created
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	tflog.Debug(ctx, "Starting Open method for Windows Autopilot Device CSV Import")

	// Create a new model to hold the configuration
	var data WindowsAutopilotDeviceCSVImportModel

	// Get the configuration from the request
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the file path
	filePath := data.FilePath.ValueString()
	if filePath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("file_path"),
			"Invalid File Path",
			"The file_path attribute cannot be empty.",
		)
		return
	}

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("file_path"),
			"Invalid File Path",
			fmt.Sprintf("Error accessing file at %s: %s", filePath, err),
		)
		return
	}

	// Check if it's a file (not a directory)
	if fileInfo.IsDir() {
		resp.Diagnostics.AddAttributeError(
			path.Root("file_path"),
			"Invalid File Path",
			fmt.Sprintf("Path %s is a directory, not a file", filePath),
		)
		return
	}

	// Read the CSV file
	devices, diags := readCSVFile(ctx, filePath)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the model with the devices
	data.Devices = devices

	// Set the result
	diags = resp.Result.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Completed Open method, read %d devices", len(devices)))
}

// Configure is called to pass the provider configured client to the resource
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Windows Autopilot Device CSV Import ephemeral resource")
	// No configuration needed for this resource
}

// ValidateConfig validates the configuration
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) ValidateConfig(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	tflog.Debug(ctx, "Validating Windows Autopilot Device CSV Import configuration")

	// Create a new model to hold the configuration
	var data WindowsAutopilotDeviceCSVImportModel

	// Get the configuration from the request
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the file path
	filePath := data.FilePath.ValueString()
	if filePath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("file_path"),
			"Invalid File Path",
			"The file_path attribute cannot be empty.",
		)
	}
}

// Close is called when the ephemeral resource is no longer needed
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	tflog.Debug(ctx, "Closing Windows Autopilot Device CSV Import ephemeral resource")
	// Nothing to clean up
}

// Renew is called when the ephemeral resource needs to be renewed
func (r *WindowsAutopilotDeviceCSVImportEphemeralResource) Renew(ctx context.Context, req ephemeral.RenewRequest, resp *ephemeral.RenewResponse) {
	tflog.Debug(ctx, "Renewing Windows Autopilot Device CSV Import ephemeral resource")
	// This resource doesn't need renewal, so we just copy the result
	// Note: Renew method doesn't have Result fields in the request/response
	// This is a no-op for this resource
}
