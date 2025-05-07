package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	// Validate that either role_definition_id or built_in_role_name is provided
	if data.RoleDefinitionID.IsNull() && data.BuiltInRoleName.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Field",
			"Either role_definition_id or built_in_role_name must be specified",
		)
		return
	}

	// Determine role definition ID
	var roleDefinitionID string
	isBuiltInRole := false
	builtInRoleName := ""

	if !data.RoleDefinitionID.IsNull() && !data.RoleDefinitionID.IsUnknown() {
		roleDefinitionID = data.RoleDefinitionID.ValueString()
	} else if !data.BuiltInRoleName.IsNull() && !data.BuiltInRoleName.IsUnknown() {
		isBuiltInRole = true
		builtInRoleName = data.BuiltInRoleName.ValueString()
	}

	// Construct the assignment
	requestBody, err := constructResource(
		ctx,
		roleDefinitionID,
		isBuiltInRole,
		builtInRoleName,
		&data,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment",
			fmt.Sprintf("Could not construct assignment: %s", err.Error()),
		)
		return
	}

	// Create the assignment via API
	createdResource, err := r.client.
		DeviceManagement().
		RoleAssignments().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create Assignment", r.WritePermissions)
		return
	}

	if createdResource.GetId() != nil {
		data.ID = types.StringValue(*createdResource.GetId())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

// Read handles the Read operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Fetch the assignment
	resource, err := r.client.
		DeviceManagement().
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(data.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Map the API response to our model
	MapRemoteResourceStateToTerraform(ctx, &data, resource)

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read method for: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s_%s", r.ProviderTypeName, r.TypeName))
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	// Determine role definition ID
	var roleDefinitionID string
	isBuiltInRole := false
	builtInRoleName := ""

	if !data.RoleDefinitionID.IsNull() && !data.RoleDefinitionID.IsUnknown() {
		roleDefinitionID = data.RoleDefinitionID.ValueString()
	} else if !data.BuiltInRoleName.IsNull() && !data.BuiltInRoleName.IsUnknown() {
		isBuiltInRole = true
		builtInRoleName = data.BuiltInRoleName.ValueString()
	}

	// Construct the updated assignment
	requestBody, err := constructResource(
		ctx,
		roleDefinitionID,
		isBuiltInRole,
		builtInRoleName,
		&data,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update",
			fmt.Sprintf("Could not construct assignment: %s", err.Error()),
		)
		return
	}

	// Update the assignment via API
	_, err = r.client.
		DeviceManagement().
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(data.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update Assignment", r.WritePermissions)
		return
	}

	// Read the updated resource to refresh state
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update method for: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleDefinitionAssignmentResourceModel

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

	// Delete the assignment
	err := r.client.
		DeviceManagement().
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete Assignment", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
