package graphbetaroledefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleDefinitionResourceModel

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

	roleDef, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	createdRoleDef, err := r.client.DeviceManagement().RoleDefinitions().Post(ctx, roleDef, nil)
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

	plan.ID = types.StringValue(*createdRoleDef.GetId())

	// Call Read to fetch the full state and set it in the response
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleDefinitionResourceModel
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

	roleDef, err := r.client.DeviceManagement().RoleDefinitions().
		ByRoleDefinitionId(state.ID.ValueString()).
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

	MapRemoteStateToTerraform(ctx, &state, roleDef)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoleDefinitionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	updatedRoleDef, err := r.client.DeviceManagement().RoleDefinitions().
		ByRoleDefinitionId(plan.ID.ValueString()).
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

	MapRemoteStateToTerraform(ctx, &plan, updatedRoleDef)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleDefinitionResourceModel

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

	err := r.client.DeviceManagement().RoleDefinitions().
		ByRoleDefinitionId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, data.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Delete", r.WritePermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
