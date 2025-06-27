package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

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

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resource, err := r.client.
		DeviceManagement().
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(data.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &data, resource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read method for: %s", ResourceName))
}

// Update handles the Update operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleDefinitionAssignmentResourceModel
	var state RoleDefinitionAssignmentResourceModel
	var roleDefinitionID string // Determine role definition ID

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	isBuiltInRole := false
	builtInRoleName := ""

	if !plan.RoleDefinitionID.IsNull() && !plan.RoleDefinitionID.IsUnknown() {
		roleDefinitionID = plan.RoleDefinitionID.ValueString()
	} else if !plan.BuiltInRoleName.IsNull() && !plan.BuiltInRoleName.IsUnknown() {
		isBuiltInRole = true
		builtInRoleName = plan.BuiltInRoleName.ValueString()
	}

	requestBody, err := constructResource(
		ctx,
		roleDefinitionID,
		isBuiltInRole,
		builtInRoleName,
		&plan,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update",
			fmt.Sprintf("Could not construct assignment: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update Assignment", r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for the RoleDefinitionAssignment resource.
func (r *RoleDefinitionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleDefinitionAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

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
		RoleAssignments().
		ByDeviceAndAppManagementRoleAssignmentId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
