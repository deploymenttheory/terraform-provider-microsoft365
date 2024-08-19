package graphbetadevicemanagementscript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/deviceandappmanagement/beta/deviceManagementScript"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *DeviceManagementScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data resource.DeviceManagementScriptResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	scriptId := data.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Device Management Script with ID: %s", scriptId))

	script, err := d.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(scriptId).Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Device Management Script, got error: %s", err))
		return
	}

	// Map the response to the data model using the resource's mapping function
	resource.MapRemoteStateToTerraform(ctx, &data, script)

	assignments, err := d.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(scriptId).Assignments().Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read Device Management Script assignments, got error: %s", err))
	} else {
		// Iterate over assignments and map each one
		for _, assignment := range assignments.GetValue() {
			resource.MapAssignmentsRemoteStateToTerraform(assignment)
		}
	}

	groupAssignments, err := d.client.DeviceManagement().DeviceManagementScripts().ByDeviceManagementScriptId(scriptId).GroupAssignments().Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read Device Management Script group assignments, got error: %s", err))
	} else {
		// Iterate over group assignments and map each one
		for _, groupAssignment := range groupAssignments.GetValue() {
			resource.MapGroupAssignmentsRemoteStateToTerraform(groupAssignment)
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
