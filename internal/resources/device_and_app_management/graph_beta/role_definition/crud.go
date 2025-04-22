package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object RoleDefinitionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, r.client, &object, resp, r.ReadPermissions, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		var assignmentsList []sharedmodels.RoleAssignmentResourceModel
		diags := object.Assignments.ElementsAs(ctx, &assignmentsList, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		for _, assignment := range assignmentsList {
			requestAssignment, err := constructAssignment(
				ctx,
				object.ID.ValueString(),
				object.IsBuiltInRoleDefinition.ValueBool(),
				object.BuiltInRoleName.ValueString(),
				&assignment,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing assignment",
					fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
				)
				return
			}

			_, err = r.client.
				DeviceManagement().
				RoleAssignments().
				Post(ctx, requestAssignment, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object RoleDefinitionResourceModel
	var assignment sharedmodels.RoleAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Get the role definition
	resource, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, resource)

	respAssignments, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		RoleAssignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read Assignments", r.ReadPermissions)
		return
	}

	MapRemoteAssignmentStateToTerraform(ctx, &assignment, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object, state RoleDefinitionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, r.client, &object, resp, r.ReadPermissions, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		var assignmentsList []sharedmodels.RoleAssignmentResourceModel
		diags := object.Assignments.ElementsAs(ctx, &assignmentsList, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		for _, assignment := range assignmentsList {
			requestAssignment, err := constructAssignment(
				ctx,
				object.ID.ValueString(),
				object.IsBuiltInRoleDefinition.ValueBool(),
				object.BuiltInRoleName.ValueString(),
				&assignment,
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing assignment",
					fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
				)
				return
			}

			_, err = r.client.
				DeviceManagement().
				RoleAssignments().
				Post(ctx, requestAssignment, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
				return
			}
		}
	}

	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State = readResp.State

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

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
