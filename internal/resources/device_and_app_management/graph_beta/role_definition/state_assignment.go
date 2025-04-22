package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"sort"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// This function strictly converts the data without making any API calls
// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// This function strictly converts the data without making any API calls
// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// This function strictly converts the data without making any API calls
func MapRemoteAssignmentStateToTerraform(ctx context.Context, roleDefinition *RoleDefinitionResourceModel, remoteAssignments graphmodels.RoleAssignmentCollectionResponseable) {
	tflog.Debug(ctx, "Starting MapRemoteAssignmentStateToTerraform")

	// Define the object type for assignments
	assignmentObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                          types.StringType,
			"display_name":                types.StringType,
			"description":                 types.StringType,
			"admin_group_users_group_ids": types.SetType{ElemType: types.StringType},
			"scope_type":                  types.StringType,
			"resource_scopes":             types.SetType{ElemType: types.StringType},
		},
	}

	// If no assignments, set to null and return early
	if remoteAssignments == nil || remoteAssignments.GetValue() == nil || len(remoteAssignments.GetValue()) == 0 {
		tflog.Debug(ctx, "No remote assignments to process, setting assignments to null")
		roleDefinition.Assignments = types.SetNull(assignmentObjectType)
		tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
		return
	}

	// Create a map to hold assignments by display name
	assignmentMap := make(map[string]sharedmodels.RoleAssignmentResourceModel)

	// Process remote assignments
	assignments := remoteAssignments.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("Processing %d remote assignments", len(assignments)))

	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		id := state.StringPtrToString(assignment.GetId())
		displayName := state.StringPtrToString(assignment.GetDisplayName())

		tflog.Debug(ctx, fmt.Sprintf("Processing assignment with ID: %s, DisplayName: %s", id, displayName))

		// Create standardized scope members and resource scopes slices (sorted alphabetically)
		var sortedScopeMembers, sortedResourceScopes []string

		// Process scope members
		if scopeMembers := assignment.GetScopeMembers(); len(scopeMembers) > 0 {
			sortedScopeMembers = make([]string, len(scopeMembers))
			copy(sortedScopeMembers, scopeMembers)
			sort.Strings(sortedScopeMembers)
			tflog.Debug(ctx, fmt.Sprintf("Assignment has %d scope members (sorted): %v",
				len(sortedScopeMembers), sortedScopeMembers))
		}

		// Process resource scopes
		if resourceScopes := assignment.GetResourceScopes(); len(resourceScopes) > 0 {
			sortedResourceScopes = make([]string, len(resourceScopes))
			copy(sortedResourceScopes, resourceScopes)
			sort.Strings(sortedResourceScopes)
			tflog.Debug(ctx, fmt.Sprintf("Assignment has %d resource scopes (sorted): %v",
				len(sortedResourceScopes), sortedResourceScopes))
		}

		// Create assignment model
		assignmentModel := sharedmodels.RoleAssignmentResourceModel{
			ID:             types.StringValue(id),
			DisplayName:    types.StringValue(displayName),
			Description:    types.StringValue(state.StringPtrToString(assignment.GetDescription())),
			ScopeMembers:   state.StringSliceToSet(ctx, sortedScopeMembers),
			ResourceScopes: state.StringSliceToSet(ctx, sortedResourceScopes),
		}

		// Process scope type (use empty string instead of null for consistency)
		if scopeType := assignment.GetScopeType(); scopeType != nil {
			assignmentModel.ScopeType = types.StringValue(string(*scopeType))
		} else {
			assignmentModel.ScopeType = types.StringValue("")
		}

		// Use display name as map key for stable ordering
		assignmentMap[displayName] = assignmentModel
	}

	// Extract assignments from map in sorted order by display name
	var displayNames []string
	for name := range assignmentMap {
		displayNames = append(displayNames, name)
	}
	sort.Strings(displayNames)

	var resultAssignments []sharedmodels.RoleAssignmentResourceModel
	for _, name := range displayNames {
		resultAssignments = append(resultAssignments, assignmentMap[name])
	}

	// Log all assignments before conversion
	for i, assignment := range resultAssignments {
		tflog.Debug(ctx, fmt.Sprintf("Assignment %d before conversion:", i), map[string]interface{}{
			"id":                          assignment.ID.ValueString(),
			"display_name":                assignment.DisplayName.ValueString(),
			"description":                 assignment.Description.ValueString(),
			"admin_group_users_group_ids": fmt.Sprintf("%v", assignment.ScopeMembers),
			"scope_type":                  assignment.ScopeType.ValueString(),
			"resource_scopes":             fmt.Sprintf("%v", assignment.ResourceScopes),
		})
	}

	// Convert assignments to Terraform set
	if len(resultAssignments) > 0 {
		// Try a more reliable approach using ObjectValueMust
		objects := make([]attr.Value, 0, len(resultAssignments))

		for _, assignment := range resultAssignments {
			obj := map[string]attr.Value{
				"id":                          assignment.ID,
				"display_name":                assignment.DisplayName,
				"description":                 assignment.Description,
				"admin_group_users_group_ids": assignment.ScopeMembers,
				"scope_type":                  assignment.ScopeType,
				"resource_scopes":             assignment.ResourceScopes,
			}

			// Create an object value
			objectValue, diags := types.ObjectValue(assignmentObjectType.AttrTypes, obj)
			if diags.HasError() {
				for _, diag := range diags.Errors() {
					tflog.Error(ctx, fmt.Sprintf("Error creating object: %s", diag.Detail()))
				}
				continue
			}

			objects = append(objects, objectValue)
		}

		// Create a set from the objects
		assignmentsSet, diags := types.SetValue(assignmentObjectType, objects)
		if diags.HasError() {
			tflog.Error(ctx, "Error converting assignments to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})

			// Set to null on error
			roleDefinition.Assignments = types.SetNull(assignmentObjectType)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Successfully converted %d assignments to set", len(resultAssignments)))
			roleDefinition.Assignments = assignmentsSet
		}
	} else {
		// No assignments, set to null
		roleDefinition.Assignments = types.SetNull(assignmentObjectType)
	}

	tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
}
