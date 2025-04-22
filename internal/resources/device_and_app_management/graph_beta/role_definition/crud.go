package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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

	// Intune roles require unique display_names
	isBuiltIn := object.IsBuiltInRoleDefinition.ValueBool() || object.IsBuiltIn.ValueBool()
	if !isBuiltIn && !object.DisplayName.IsNull() && !object.DisplayName.IsUnknown() {
		if err := checkRoleNameUniqueness(ctx, r.client, object.DisplayName.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Role Name Not Unique",
				err.Error(),
			)
			return
		}
	}

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

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// 1️⃣ Fetch base resource
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

	// 2️⃣ List the assignments
	assignmentsList, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		RoleAssignments().
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read Assignments List", r.ReadPermissions)
		return
	}

	// 3️⃣ Pull each assignment’s full details
	detailedResponse := graphmodels.NewRoleAssignmentCollectionResponse()
	var detailedAssignments []graphmodels.RoleAssignmentable

	if assignmentsList != nil && assignmentsList.GetValue() != nil {
		for _, listAssignment := range assignmentsList.GetValue() {
			if listAssignment == nil || listAssignment.GetId() == nil {
				continue
			}
			assignmentID := *listAssignment.GetId()
			full, err := r.client.
				DeviceManagement().
				RoleAssignments().
				ByDeviceAndAppManagementRoleAssignmentId(assignmentID).
				Get(ctx, nil)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to fetch details for assignment ID %s: %s", assignmentID, err))
				continue
			}
			// append to the RoleAssignmentable slice
			detailedAssignments = append(detailedAssignments, full)
		}
	}
	detailedResponse.SetValue(detailedAssignments)

	// ─── DEBUG DUMP ────────────────────────────────────────────────────────────────
	debugPrintAssignments(ctx, detailedAssignments)
	// ───────────────────────────────────────────────────────────────────────────────

	// 4️⃣ Finally map into Terraform state
	MapRemoteAssignmentStateToTerraform(ctx, &object, detailedResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read method for: %s_%s", r.ProviderTypeName, r.TypeName))
}

// debugPrintAssignments emits a detailed tflog.Debug of every assignment object
func debugPrintAssignments(ctx context.Context, assigns []graphmodels.RoleAssignmentable) {
	tflog.Debug(ctx, fmt.Sprintf("=== Graph API returned %d assignments:", len(assigns)))
	for i, a := range assigns {
		id := state.StringPtrToString(a.GetId())
		name := state.StringPtrToString(a.GetDisplayName())
		desc := state.StringPtrToString(a.GetDescription())
		members := a.GetScopeMembers()
		scopes := a.GetResourceScopes()
		scopeType := "<nil>"
		if st := a.GetScopeType(); st != nil {
			scopeType = string(*st)
		}
		tflog.Debug(ctx, fmt.Sprintf(
			"Assignment[%d]: id=%q, displayName=%q, description=%q, scopeType=%q, scopeMembers=%v, resourceScopes=%v",
			i, id, name, desc, scopeType, members, scopes,
		))
	}
	tflog.Debug(ctx, "=== end of API assignment dump")
}

// Update handles the Update operation for role definitions and assignments,
// tracking assignments strictly by ID
// Update handles the Update operation for role definitions and assignments,
// tracking assignments strictly by ID and tracking IDs for new assignments
func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planObj, stateObj RoleDefinitionResourceModel

	tflog.Info(ctx, "Starting Update for RoleDefinition", map[string]interface{}{"resource_type": r.TypeName})

	// Load plan & state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planObj)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateObj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, planObj.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	// 1️⃣ PATCH the base RoleDefinition
	builder := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(planObj.ID.ValueString())
	requestBody, err := constructResource(ctx, r.client, &planObj, resp, r.ReadPermissions, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource", err.Error())
		return
	}
	if _, err := builder.Patch(ctx, requestBody, nil); err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update RoleDefinition", r.WritePermissions)
		return
	}
	tflog.Debug(ctx, "Patched base RoleDefinition successfully")

	// 2️⃣ Extract assignments from state and plan
	var stateAssignments, planAssignments []sharedmodels.RoleAssignmentResourceModel

	if !stateObj.Assignments.IsNull() && !stateObj.Assignments.IsUnknown() {
		diags := stateObj.Assignments.ElementsAs(ctx, &stateAssignments, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !planObj.Assignments.IsNull() && !planObj.Assignments.IsUnknown() {
		diags := planObj.Assignments.ElementsAs(ctx, &planAssignments, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// 3️⃣ Create maps for state by ID
	stateAssignByID := make(map[string]sharedmodels.RoleAssignmentResourceModel)
	stateIDs := make(map[string]bool) // Track all state IDs

	// Track all assignments in state by ID
	for _, a := range stateAssignments {
		if !a.ID.IsNull() && !a.ID.IsUnknown() {
			id := a.ID.ValueString()
			stateAssignByID[id] = a
			stateIDs[id] = true
		}
	}

	// 4️⃣ Process plan assignments into operations
	toUpdate := make(map[string]sharedmodels.RoleAssignmentResourceModel) // ID -> assignment to update
	toCreate := make([]sharedmodels.RoleAssignmentResourceModel, 0)       // Assignments to create
	planIDs := make(map[string]bool)                                      // Track which IDs are in plan

	// Determine which plan assignments need update vs create
	for _, planAssign := range planAssignments {
		if !planAssign.ID.IsNull() && !planAssign.ID.IsUnknown() {
			// This plan assignment has an ID - track it for update
			id := planAssign.ID.ValueString()
			toUpdate[id] = planAssign
			planIDs[id] = true
		} else {
			// No ID - this is a new assignment
			toCreate = append(toCreate, planAssign)
		}
	}

	// 5️⃣ Determine which state assignments to delete
	toDelete := make([]string, 0)
	for id := range stateIDs {
		if !planIDs[id] {
			// This ID exists in state but not in plan
			toDelete = append(toDelete, id)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Assignment operations: %d to create, %d to update, %d to delete",
		len(toCreate), len(toUpdate), len(toDelete)))

	// 6️⃣ Delete assignments that should be removed
	for _, id := range toDelete {
		tflog.Debug(ctx, "Deleting assignment", map[string]interface{}{"id": id})
		err := r.client.
			DeviceManagement().
			RoleAssignments().
			ByDeviceAndAppManagementRoleAssignmentId(id).
			Delete(ctx, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, fmt.Sprintf("Delete Assignment %s", id), r.WritePermissions)
			return
		}
	}

	// 7️⃣ Update existing assignments
	for id, planAssign := range toUpdate {
		displayName := "<unknown>"
		if !planAssign.DisplayName.IsNull() && !planAssign.DisplayName.IsUnknown() {
			displayName = planAssign.DisplayName.ValueString()
		}

		tflog.Debug(ctx, "Updating assignment", map[string]interface{}{
			"id":           id,
			"display_name": displayName,
		})

		updatedAssign, err := constructAssignment(
			ctx,
			planObj.ID.ValueString(),
			planObj.IsBuiltInRoleDefinition.ValueBool(),
			planObj.BuiltInRoleName.ValueString(),
			&planAssign,
		)
		if err != nil {
			resp.Diagnostics.AddError("Error constructing assignment for update", err.Error())
			return
		}

		_, err = r.client.
			DeviceManagement().
			RoleAssignments().
			ByDeviceAndAppManagementRoleAssignmentId(id).
			Patch(ctx, updatedAssign, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, fmt.Sprintf("Update Assignment %s", id), r.WritePermissions)
			return
		}
	}

	// 8️⃣ Create new assignments and collect their IDs
	newAssignmentIDs := make(map[string]string) // Map display names to created IDs

	for _, planAssign := range toCreate {
		displayName := "<unknown>"
		if !planAssign.DisplayName.IsNull() && !planAssign.DisplayName.IsUnknown() {
			displayName = planAssign.DisplayName.ValueString()
		}

		tflog.Debug(ctx, "Creating new assignment", map[string]interface{}{
			"display_name": displayName,
		})

		// Ensure ID is null for new assignments
		newAssign := planAssign
		newAssign.ID = types.StringNull()

		reqBody, err := constructAssignment(
			ctx,
			planObj.ID.ValueString(),
			planObj.IsBuiltInRoleDefinition.ValueBool(),
			planObj.BuiltInRoleName.ValueString(),
			&newAssign,
		)
		if err != nil {
			resp.Diagnostics.AddError("Error constructing assignment for create", err.Error())
			return
		}

		// Create the assignment
		createdAssign, err := r.client.
			DeviceManagement().
			RoleAssignments().
			Post(ctx, reqBody, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create Assignment", r.WritePermissions)
			return
		}

		// Store the new ID for state tracking
		if createdAssign != nil && createdAssign.GetId() != nil {
			newID := *createdAssign.GetId()
			newAssignmentIDs[displayName] = newID
			tflog.Debug(ctx, "Created assignment with ID", map[string]interface{}{
				"display_name": displayName,
				"id":           newID,
			})
		}
	}

	// 9️⃣ Let Read function handle state
	tflog.Debug(ctx, "Using Read to refresh final state")
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = readResp.State

	tflog.Info(ctx, "Update completed successfully")
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

	isBuiltIn := data.IsBuiltInRoleDefinition.ValueBool() || data.IsBuiltIn.ValueBool()

	// For built-in roles, we only need to delete the assignments
	if isBuiltIn {
		tflog.Debug(ctx, "Built-in role detected - will only delete assignments")

		respAssignments, err := r.client.
			DeviceManagement().
			RoleDefinitions().
			ByRoleDefinitionId(data.ID.ValueString()).
			RoleAssignments().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Get Assignments", r.ReadPermissions)
			return
		}

		assignments := respAssignments.GetValue()
		for _, assignment := range assignments {
			assignmentID := *assignment.GetId()
			tflog.Debug(ctx, fmt.Sprintf("Deleting assignment with ID: %s", assignmentID))

			err := r.client.
				DeviceManagement().
				RoleAssignments().
				ByDeviceAndAppManagementRoleAssignmentId(assignmentID).
				Delete(ctx, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, fmt.Sprintf("Delete Assignment %s", assignmentID), r.WritePermissions)
				return
			}
		}

		tflog.Debug(ctx, "All assignments for built-in role deleted successfully")
	} else {
		tflog.Debug(ctx, "Custom role detected - will delete the entire role definition")

		err := r.client.
			DeviceManagement().
			RoleDefinitions().
			ByRoleDefinitionId(data.ID.ValueString()).
			Delete(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Custom role definition deleted successfully")
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
