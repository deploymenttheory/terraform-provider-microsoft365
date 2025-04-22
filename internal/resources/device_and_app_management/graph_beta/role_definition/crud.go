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

// Update handles the Update operation for the RoleDefinition resource, performing differential patch/create/delete of assignments.
func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planObj, stateObj RoleDefinitionResourceModel

	// Add logging about the update process starting
	tflog.Info(ctx, "Starting Update for RoleDefinition", map[string]interface{}{
		"resource_type": r.TypeName,
	})

	// Load both plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planObj)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateObj)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to extract plan or state objects", map[string]interface{}{
			"has_errors": resp.Diagnostics.HasError(),
		})
		return
	}

	tflog.Debug(ctx, "Plan and state loaded successfully", map[string]interface{}{
		"plan_id":  planObj.ID.ValueString(),
		"state_id": stateObj.ID.ValueString(),
	})

	// 1️⃣ Patch the base RoleDefinition
	requestBody, err := constructResource(ctx, r.client, &planObj, resp, r.ReadPermissions, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource", err.Error())
		return
	}
	_, err = r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(planObj.ID.ValueString()).
		Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully patched base role definition")

	// 2️⃣ Fetch existing assignments from the API
	tflog.Debug(ctx, "Fetching existing role assignments")
	listResp, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(planObj.ID.ValueString()).
		RoleAssignments().
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Get Assignments", r.ReadPermissions)
		return
	}

	// Map existing assignments by ID for easier lookup
	existingByID := make(map[string]graphmodels.RoleAssignmentable, len(listResp.GetValue()))
	for _, a := range listResp.GetValue() {
		if id := a.GetId(); id != nil {
			existingByID[*id] = a
			tflog.Debug(ctx, "Found existing assignment", map[string]interface{}{
				"id":           *id,
				"display_name": state.StringPtrToString(a.GetDisplayName()),
			})
		}
	}
	tflog.Debug(ctx, "Finished mapping existing assignments", map[string]interface{}{
		"count": len(existingByID),
	})

	// 3️⃣ Pull desired assignments out of the plan
	var desiredList []sharedmodels.RoleAssignmentResourceModel
	if !planObj.Assignments.IsNull() && !planObj.Assignments.IsUnknown() {
		if diags := planObj.Assignments.ElementsAs(ctx, &desiredList, false); diags.HasError() {
			tflog.Error(ctx, "Failed to extract assignments from plan", map[string]interface{}{
				"error_count": len(diags.Errors()),
			})
			resp.Diagnostics.Append(diags...)
			return
		}
		tflog.Debug(ctx, "Extracted desired assignments from plan", map[string]interface{}{
			"count": len(desiredList),
		})
	} else {
		tflog.Debug(ctx, "No assignments found in plan object")
	}

	// Log the full desiredList for debugging
	for i, assignment := range desiredList {
		tflog.Debug(ctx, "Desired assignment details", map[string]interface{}{
			"index":        i,
			"id":           assignment.ID.String(),
			"display_name": assignment.DisplayName.ValueString(),
			"description":  assignment.Description.ValueString(),
			"has_id":       !assignment.ID.IsNull() && !assignment.ID.IsUnknown(),
			"id_null":      assignment.ID.IsNull(),
			"id_unknown":   assignment.ID.IsUnknown(),
		})
	}

	// Categorize assignments into updates, deletes, and creates
	toUpdateByID := make(map[string]sharedmodels.RoleAssignmentResourceModel)
	var toCreate []sharedmodels.RoleAssignmentResourceModel
	var toDelete []string

	// Find assignments to update or create
	for i, d := range desiredList {
		if !d.ID.IsNull() && !d.ID.IsUnknown() {
			// Check if this ID is already in the map
			if existing, exists := toUpdateByID[d.ID.ValueString()]; exists {
				tflog.Warn(ctx, "Duplicate assignment ID detected", map[string]interface{}{
					"id":                    d.ID.ValueString(),
					"existing_display_name": existing.DisplayName.ValueString(),
					"new_display_name":      d.DisplayName.ValueString(),
					"index":                 i,
				})

				// If the ID is the same but display name is different, treat it as a new assignment to create
				if existing.DisplayName.ValueString() != d.DisplayName.ValueString() {
					// Create a new assignment without the ID
					newAssignment := d
					newAssignment.ID = types.StringNull()
					toCreate = append(toCreate, newAssignment)
					tflog.Info(ctx, "Moving assignment with duplicate ID to creation list", map[string]interface{}{
						"display_name": d.DisplayName.ValueString(),
					})
					continue
				}
			}

			// If this ID exists in the API, it's an update; otherwise, it's likely an error or a rename
			if _, exists := existingByID[d.ID.ValueString()]; exists {
				toUpdateByID[d.ID.ValueString()] = d
				tflog.Debug(ctx, "Assignment with ID marked for update", map[string]interface{}{
					"index":        i,
					"id":           d.ID.ValueString(),
					"display_name": d.DisplayName.ValueString(),
				})
			} else {
				// ID in plan doesn't exist in API - could be an error or a rename
				// To be safe, treat it as a new creation
				newAssignment := d
				newAssignment.ID = types.StringNull()
				toCreate = append(toCreate, newAssignment)
				tflog.Warn(ctx, "Assignment ID in plan not found in API, treating as new", map[string]interface{}{
					"id":           d.ID.ValueString(),
					"display_name": d.DisplayName.ValueString(),
				})
			}
		} else {
			// No ID provided, definitely a new creation
			toCreate = append(toCreate, d)
			tflog.Debug(ctx, "Assignment without ID marked for creation", map[string]interface{}{
				"index":        i,
				"display_name": d.DisplayName.ValueString(),
			})
		}
	}

	// Find assignments to delete (in API but not in plan)
	for id := range existingByID {
		if _, keep := toUpdateByID[id]; !keep {
			toDelete = append(toDelete, id)
			tflog.Debug(ctx, "Assignment marked for deletion", map[string]interface{}{
				"id":           id,
				"display_name": state.StringPtrToString(existingByID[id].GetDisplayName()),
			})
		}
	}

	tflog.Debug(ctx, "Assignments categorized", map[string]interface{}{
		"to_update_count": len(toUpdateByID),
		"to_create_count": len(toCreate),
		"to_delete_count": len(toDelete),
	})

	// STEP 1: PATCH existing assignments
	tflog.Info(ctx, "Processing assignments to update", map[string]interface{}{
		"count": len(toUpdateByID),
	})

	for id, desired := range toUpdateByID {
		tflog.Debug(ctx, "Patching assignment", map[string]interface{}{
			"id":           id,
			"display_name": desired.DisplayName.ValueString(),
		})
		reqBody, err := constructAssignment(
			ctx,
			planObj.ID.ValueString(),
			planObj.IsBuiltInRoleDefinition.ValueBool(),
			planObj.BuiltInRoleName.ValueString(),
			&desired,
		)
		if err != nil {
			resp.Diagnostics.AddError("Error constructing assignment", err.Error())
			return
		}
		if _, err := r.client.
			DeviceManagement().
			RoleAssignments().
			ByDeviceAndAppManagementRoleAssignmentId(id).
			Patch(ctx, reqBody, nil); err != nil {
			errors.HandleGraphError(ctx, err, resp, fmt.Sprintf("Patch Assignment %s", id), r.WritePermissions)
			return
		}
		tflog.Debug(ctx, "Successfully patched assignment", map[string]interface{}{
			"id": id,
		})
	}

	// Refresh state after updates
	if len(toUpdateByID) > 0 {
		tflog.Debug(ctx, "Refreshing state after updates")
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
		resp.Diagnostics.Append(readResp.Diagnostics...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Errors encountered during state refresh after updates", map[string]interface{}{
				"error_count": len(readResp.Diagnostics.Errors()),
			})
			return
		}
		resp.State = readResp.State
	}

	// STEP 2: DELETE assignments
	tflog.Info(ctx, "Processing assignments to delete", map[string]interface{}{
		"count": len(toDelete),
	})

	for _, id := range toDelete {
		tflog.Debug(ctx, "Deleting assignment", map[string]interface{}{
			"id":           id,
			"display_name": state.StringPtrToString(existingByID[id].GetDisplayName()),
		})
		if err := r.client.
			DeviceManagement().
			RoleAssignments().
			ByDeviceAndAppManagementRoleAssignmentId(id).
			Delete(ctx, nil); err != nil {
			errors.HandleGraphError(ctx, err, resp, fmt.Sprintf("Delete Assignment %s", id), r.WritePermissions)
			return
		}
		tflog.Debug(ctx, "Successfully deleted assignment", map[string]interface{}{
			"id": id,
		})
	}

	// Refresh state after deletes
	if len(toDelete) > 0 {
		tflog.Debug(ctx, "Refreshing state after deletes")
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
		resp.Diagnostics.Append(readResp.Diagnostics...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Errors encountered during state refresh after deletes", map[string]interface{}{
				"error_count": len(readResp.Diagnostics.Errors()),
			})
			return
		}
		resp.State = readResp.State
	}

	// STEP 3: CREATE new assignments
	tflog.Info(ctx, "Processing assignments to create", map[string]interface{}{
		"count": len(toCreate),
	})

	for i, newA := range toCreate {
		tflog.Debug(ctx, "Creating new assignment", map[string]interface{}{
			"index":               i,
			"display_name":        newA.DisplayName.ValueString(),
			"scope_type":          newA.ScopeType.ValueString(),
			"has_scope_members":   !newA.ScopeMembers.IsNull() && !newA.ScopeMembers.IsUnknown(),
			"has_resource_scopes": !newA.ResourceScopes.IsNull() && !newA.ResourceScopes.IsUnknown(),
		})

		// Dump scope members and resource scopes for detailed debugging
		if !newA.ScopeMembers.IsNull() && !newA.ScopeMembers.IsUnknown() {
			var members []string
			newA.ScopeMembers.ElementsAs(ctx, &members, false)
			tflog.Debug(ctx, "Scope members", map[string]interface{}{
				"count":   len(members),
				"members": members,
			})
		}

		if !newA.ResourceScopes.IsNull() && !newA.ResourceScopes.IsUnknown() {
			var scopes []string
			newA.ResourceScopes.ElementsAs(ctx, &scopes, false)
			tflog.Debug(ctx, "Resource scopes", map[string]interface{}{
				"count":  len(scopes),
				"scopes": scopes,
			})
		}

		reqBody, err := constructAssignment(
			ctx,
			planObj.ID.ValueString(),
			planObj.IsBuiltInRoleDefinition.ValueBool(),
			planObj.BuiltInRoleName.ValueString(),
			&newA,
		)
		if err != nil {
			tflog.Error(ctx, "Failed to construct assignment", map[string]interface{}{
				"error":        err.Error(),
				"display_name": newA.DisplayName.ValueString(),
			})
			resp.Diagnostics.AddError("Error constructing assignment", err.Error())
			return
		}

		tflog.Debug(ctx, "Sending POST request to create new assignment", map[string]interface{}{
			"role_def_id":  planObj.ID.ValueString(),
			"display_name": newA.DisplayName.ValueString(),
		})

		createdAssignment, err := r.client.
			DeviceManagement().
			RoleAssignments().
			Post(ctx, reqBody, nil)

		if err != nil {
			tflog.Error(ctx, "Failed to create assignment", map[string]interface{}{
				"error":        err.Error(),
				"display_name": newA.DisplayName.ValueString(),
			})
			errors.HandleGraphError(ctx, err, resp, "Create Assignment", r.WritePermissions)
			return
		}

		if createdAssignment != nil && createdAssignment.GetId() != nil {
			tflog.Info(ctx, "Successfully created new assignment", map[string]interface{}{
				"display_name": newA.DisplayName.ValueString(),
				"new_id":       *createdAssignment.GetId(),
			})
		} else {
			tflog.Warn(ctx, "Created assignment but ID is nil or missing", map[string]interface{}{
				"display_name": newA.DisplayName.ValueString(),
			})
		}

		// Refresh state after each create to ensure the new ID is properly captured
		tflog.Debug(ctx, "Refreshing state after creating assignment", map[string]interface{}{
			"display_name": newA.DisplayName.ValueString(),
		})
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
		resp.Diagnostics.Append(readResp.Diagnostics...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Errors encountered during state refresh after create", map[string]interface{}{
				"error_count": len(readResp.Diagnostics.Errors()),
			})
			return
		}
		resp.State = readResp.State
	}

	// Final state refresh to ensure everything is consistent
	tflog.Debug(ctx, "Performing final state refresh")
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Errors encountered during final state refresh", map[string]interface{}{
			"error_count": len(readResp.Diagnostics.Errors()),
		})
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
