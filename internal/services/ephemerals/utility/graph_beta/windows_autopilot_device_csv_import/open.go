package windowsAutopilotDeviceCSVImport

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

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
	// No configuration needed for this resource as it only reads from CSV files
	tflog.Debug(ctx, "Configure called for Windows Autopilot Device CSV Import ephemeral resource (no-op)")
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
