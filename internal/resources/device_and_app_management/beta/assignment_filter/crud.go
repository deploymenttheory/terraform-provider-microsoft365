package graphBetaAssignmentFilter

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *AssignmentFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	resource, err := r.client.DeviceManagement().AssignmentFilters().Post(ctx, requestBody, nil)
	if err != nil {
		if crud.PermissionError(err, "Create", r.WritePermissions, resp) {
			return
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName),
				err.Error(),
			)
		}
		return
	}

	plan.ID = types.StringValue(*resource.GetId())

	MapRemoteStateToTerraform(ctx, &plan, resource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *AssignmentFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssignmentFilterResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resource, err := r.client.DeviceManagement().AssignmentFilters().
		ByDeviceAndAppManagementAssignmentFilterId(state.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Read", r.ReadPermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, resource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *AssignmentFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceManagement().AssignmentFilters().
		ByDeviceAndAppManagementAssignmentFilterId(plan.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, plan.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Update", r.WritePermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *AssignmentFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
