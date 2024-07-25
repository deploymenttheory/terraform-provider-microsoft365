package graphBetaAssignmentFilter

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *AssignmentFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 30*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	requestBody, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment filter",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	assignmentFilter, err := r.client.DeviceManagement().AssignmentFilters().Post(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment filter",
			fmt.Sprintf("Could not create assignment filter: %s", err.Error()),
		)
		return
	}

	data.ID = types.StringValue(*assignmentFilter.GetId())

	r.isCreate = true

	readResp := resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)

	r.isCreate = false

	tflog.Debug(ctx, fmt.Sprintf("Finished creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the read operation and stating.
func (r *AssignmentFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AssignmentFilterResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, 30*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	remoteResource, err := r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		if common.IsNotFoundError(err) && !r.isCreate {
			resp.Diagnostics.AddWarning(
				"Resource Not Found",
				fmt.Sprintf("The resource: %s_%s with ID %s was not found and will be removed from the state.", r.ProviderTypeName, r.TypeName, data.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading assignment filter",
			fmt.Sprintf("Could not read resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	mapRemoteStateToTerraform(ctx, &data, remoteResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles the Update operation.
func (r *AssignmentFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := data.Timeouts.Update(ctx, 30*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	requestBody, err := constructResource(ctx, &data)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment filter",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		if common.IsNotFoundError(err) && !r.isCreate {
			resp.Diagnostics.AddWarning(
				"Resource Not Found",
				fmt.Sprintf("The resource: %s_%s with ID %s was not found and will be removed from the state.", r.ProviderTypeName, r.TypeName, data.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading assignment filter",
			fmt.Sprintf("Could not update resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *AssignmentFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, 30*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	err := r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Completed deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
