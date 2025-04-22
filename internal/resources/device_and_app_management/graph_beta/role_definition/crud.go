package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"sort"
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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Get the base role definition
	resource, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Map the base resource state
	MapRemoteResourceStateToTerraform(ctx, &object, resource)

	// Get the list of assignments for this role definition
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

	// Create a collection of detailed assignment objects to pass to MapRemoteAssignmentStateToTerraform
	detailedResponse := graphmodels.NewRoleAssignmentCollectionResponse()
	var detailedAssignments []graphmodels.RoleAssignmentable

	// Fetch each assignment individually to get complete details
	if assignmentsList != nil && assignmentsList.GetValue() != nil {
		for _, listAssignment := range assignmentsList.GetValue() {
			if listAssignment == nil || listAssignment.GetId() == nil {
				continue
			}

			assignmentID := *listAssignment.GetId()
			tflog.Debug(ctx, fmt.Sprintf("Fetching details for assignment ID: %s", assignmentID))

			// Get the full assignment details for this ID
			detailedAssignment, err := r.client.
				DeviceManagement().
				RoleAssignments().
				ByDeviceAndAppManagementRoleAssignmentId(assignmentID).
				Get(ctx, nil)

			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to fetch details for assignment ID %s: %s", assignmentID, err.Error()))
				continue
			}

			// Log the detailed data to compare with list data
			tflog.Debug(ctx, fmt.Sprintf("Detailed assignment data for ID %s:", assignmentID), map[string]interface{}{
				"displayName":    state.StringPtrToString(detailedAssignment.GetDisplayName()),
				"scopeMembers":   detailedAssignment.GetScopeMembers(),
				"resourceScopes": detailedAssignment.GetResourceScopes(),
				"hasScopeType":   detailedAssignment.GetScopeType() != nil,
			})

			// Add to our collection
			detailedAssignments = append(detailedAssignments, detailedAssignment)
		}
	}

	// Set the values in the response
	detailedResponse.SetValue(detailedAssignments)

	// Map the detailed assignments to the Terraform state
	MapRemoteAssignmentStateToTerraform(ctx, &object, detailedResponse)

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

	// Update the base role definition
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

	// Get list of existing assignments
	respAssignments, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		RoleAssignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Get Assignments", r.ReadPermissions)
		return
	}

	// Delete all existing assignments
	assignments := respAssignments.GetValue()
	for _, assignment := range assignments {
		if assignment == nil || assignment.GetId() == nil {
			continue
		}

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

	// Get new assignments from plan
	var newAssignments []sharedmodels.RoleAssignmentResourceModel
	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		diags := object.Assignments.ElementsAs(ctx, &newAssignments, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Sort assignments for consistency
		sort.Slice(newAssignments, func(i, j int) bool {
			return newAssignments[i].DisplayName.ValueString() < newAssignments[j].DisplayName.ValueString()
		})
	}

	// Create all new assignments
	for _, assignment := range newAssignments {
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

		// Always create as new
		_, err = r.client.
			DeviceManagement().
			RoleAssignments().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create Assignment", r.WritePermissions)
			return
		}
	}

	// Allow some time for the backend to process the assignments
	time.Sleep(2 * time.Second)

	// Read back the full state
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

// func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var object, state RoleDefinitionResourceModel

// 	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
// 	if cancel == nil {
// 		return
// 	}
// 	defer cancel()

// 	requestBody, err := constructResource(ctx, r.client, &object, resp, r.ReadPermissions, true)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error constructing resource",
// 			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
// 		)
// 		return
// 	}

// 	_, err = r.client.
// 		DeviceManagement().
// 		RoleDefinitions().
// 		ByRoleDefinitionId(object.ID.ValueString()).
// 		Patch(ctx, requestBody, nil)

// 	if err != nil {
// 		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
// 		return
// 	}

// 	// Get existing assignments from state
// 	var existingAssignments []sharedmodels.RoleAssignmentResourceModel
// 	if !state.Assignments.IsNull() && !state.Assignments.IsUnknown() {
// 		diags := state.Assignments.ElementsAs(ctx, &existingAssignments, false)
// 		if diags.HasError() {
// 			resp.Diagnostics.Append(diags...)
// 			return
// 		}
// 	}

// 	existingAssignmentMap := make(map[string]sharedmodels.RoleAssignmentResourceModel)
// 	for _, assignment := range existingAssignments {
// 		if !assignment.ID.IsNull() && !assignment.ID.IsUnknown() {
// 			existingAssignmentMap[assignment.ID.ValueString()] = assignment
// 		}
// 	}

// 	// Get new assignments from plan
// 	var newAssignments []sharedmodels.RoleAssignmentResourceModel
// 	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
// 		diags := object.Assignments.ElementsAs(ctx, &newAssignments, false)
// 		if diags.HasError() {
// 			resp.Diagnostics.Append(diags...)
// 			return
// 		}
// 	}

// 	newAssignmentMap := make(map[string]bool)
// 	for _, assignment := range newAssignments {
// 		if !assignment.DisplayName.IsNull() && !assignment.DisplayName.IsUnknown() {
// 			newAssignmentMap[assignment.DisplayName.ValueString()] = true
// 		}
// 	}

// 	// Find assignments to delete (in existing but not in new)
// 	for assignmentId, existingAssignment := range existingAssignmentMap {
// 		if !existingAssignment.DisplayName.IsNull() && !existingAssignment.DisplayName.IsUnknown() {
// 			displayName := existingAssignment.DisplayName.ValueString()
// 			if _, exists := newAssignmentMap[displayName]; !exists {
// 				tflog.Debug(ctx, fmt.Sprintf("Deleting assignment with ID: %s, DisplayName: %s", assignmentId, displayName))

// 				err := r.client.
// 					DeviceManagement().
// 					RoleAssignments().
// 					ByDeviceAndAppManagementRoleAssignmentId(assignmentId).
// 					Delete(ctx, nil)

// 				if err != nil {
// 					errors.HandleGraphError(ctx, err, resp, "Delete Assignment", r.WritePermissions)
// 					return
// 				}
// 			}
// 		}
// 	}

// 	// Create or update assignments
// 	for _, assignment := range newAssignments {
// 		requestAssignment, err := constructAssignment(
// 			ctx,
// 			object.ID.ValueString(),
// 			object.IsBuiltInRoleDefinition.ValueBool(),
// 			object.BuiltInRoleName.ValueString(),
// 			&assignment,
// 		)
// 		if err != nil {
// 			resp.Diagnostics.AddError(
// 				"Error constructing assignment",
// 				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
// 			)
// 			return
// 		}

// 		if !assignment.ID.IsNull() && !assignment.ID.IsUnknown() {
// 			_, err = r.client.
// 				DeviceManagement().
// 				RoleAssignments().
// 				ByDeviceAndAppManagementRoleAssignmentId(assignment.ID.ValueString()).
// 				Patch(ctx, requestAssignment, nil)
// 		} else {
// 			_, err = r.client.
// 				DeviceManagement().
// 				RoleAssignments().
// 				Post(ctx, requestAssignment, nil)
// 		}

// 		if err != nil {
// 			operation := "Create"
// 			if !assignment.ID.IsNull() && !assignment.ID.IsUnknown() {
// 				operation = "Update"
// 			}
// 			errors.HandleGraphError(ctx, err, resp, operation+" Assignment", r.WritePermissions)
// 			return
// 		}
// 	}

// 	readResp := &resource.ReadResponse{
// 		State: resp.State,
// 	}
// 	r.Read(ctx, resource.ReadRequest{
// 		State:        resp.State,
// 		ProviderMeta: req.ProviderMeta,
// 	}, readResp)

// 	resp.Diagnostics.Append(readResp.Diagnostics...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resp.State = readResp.State

// 	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
// }

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
