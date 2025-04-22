package graphBetaRoleDefinition

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// while preserving the original plan/config data during creation and updates
func MapRemoteAssignmentStateToTerraform(ctx context.Context, roleDefinition *RoleDefinitionResourceModel, remoteAssignments graphmodels.RoleAssignmentCollectionResponseable) {
	tflog.Debug(ctx, "Starting MapRemoteAssignmentStateToTerraform")

	// Extract planned assignments if they exist
	plannedAssignments := make(map[string]sharedmodels.RoleAssignmentResourceModel)
	var planAssignmentsList []sharedmodels.RoleAssignmentResourceModel

	if !roleDefinition.Assignments.IsNull() && !roleDefinition.Assignments.IsUnknown() {
		diags := roleDefinition.Assignments.ElementsAs(ctx, &planAssignmentsList, false)
		if !diags.HasError() {
			// Create a map of assignments by display name for lookup
			for _, assignment := range planAssignmentsList {
				if !assignment.DisplayName.IsNull() && !assignment.DisplayName.IsUnknown() {
					plannedAssignments[assignment.DisplayName.ValueString()] = assignment
				}
			}
			tflog.Debug(ctx, "Extracted planned assignments", map[string]interface{}{
				"count": len(plannedAssignments),
			})
		}
	}

	// Create a map to track assignments that have been processed
	processedAssignments := make(map[string]bool)

	// Process remote assignments if available
	var resultAssignments []sharedmodels.RoleAssignmentResourceModel

	if remoteAssignments != nil {
		assignments := remoteAssignments.GetValue()
		if len(assignments) > 0 {
			tflog.Debug(ctx, "Processing remote assignments", map[string]interface{}{
				"count": len(assignments),
			})

			for _, assignment := range assignments {
				if assignment == nil {
					continue
				}

				// Get basic properties from the remote assignment
				assignmentID := state.StringPtrToString(assignment.GetId())
				displayName := state.StringPtrToString(assignment.GetDisplayName())

				// Create a new model - we'll either merge with planned data or use remote data
				resultAssignment := sharedmodels.RoleAssignmentResourceModel{
					ID: types.StringValue(assignmentID),
				}

				// Check if we have this assignment in the plan (by display name)
				if plannedAssignment, exists := plannedAssignments[displayName]; exists {
					tflog.Debug(ctx, "Found matching planned assignment", map[string]interface{}{
						"displayName": displayName,
						"id":          assignmentID,
					})

					// Start with the planned assignment to preserve user-provided values
					resultAssignment = plannedAssignment

					// Update only the ID from the remote assignment
					resultAssignment.ID = types.StringValue(assignmentID)

					// Mark this assignment as processed
					processedAssignments[displayName] = true
				} else {
					// This is a remote assignment that wasn't in the plan
					// Likely created outside Terraform or a rename occurred
					tflog.Debug(ctx, "Remote assignment not in plan", map[string]interface{}{
						"displayName": displayName,
						"id":          assignmentID,
					})

					// Set basic properties from remote assignment
					resultAssignment.DisplayName = types.StringValue(displayName)
					resultAssignment.Description = types.StringValue(state.StringPtrToString(assignment.GetDescription()))

					// Map members to admin_group_users_group_ids (ScopeMembers in struct)
					members := assignment.GetScopeMembers()
					if len(members) > 0 {
						membersSet, diags := types.SetValueFrom(ctx, types.StringType, members)
						if !diags.HasError() {
							resultAssignment.ScopeMembers = membersSet
						} else {
							emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
							resultAssignment.ScopeMembers = emptySet
						}
					} else {
						emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
						resultAssignment.ScopeMembers = emptySet
					}

					// Map resource scopes
					resourceScopes := assignment.GetResourceScopes()
					if len(resourceScopes) > 0 {
						resourceScopesSet, diags := types.SetValueFrom(ctx, types.StringType, resourceScopes)
						if !diags.HasError() {
							resultAssignment.ResourceScopes = resourceScopesSet
						} else {
							emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
							resultAssignment.ResourceScopes = emptySet
						}
					} else {
						emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
						resultAssignment.ResourceScopes = emptySet
					}

					// Map scope type
					if scopeType := assignment.GetScopeType(); scopeType != nil {
						resultAssignment.ScopeType = types.StringValue(scopeType.String())
					} else {
						resultAssignment.ScopeType = types.StringNull()
					}
				}

				// Add to result list
				resultAssignments = append(resultAssignments, resultAssignment)
			}
		}
	}

	// Add planned assignments that weren't found in remote assignments
	// This is important for newly created assignments that might not be returned yet
	for displayName, plannedAssignment := range plannedAssignments {
		if !processedAssignments[displayName] {
			tflog.Debug(ctx, "Adding planned assignment not found in remote data", map[string]interface{}{
				"displayName": displayName,
			})
			resultAssignments = append(resultAssignments, plannedAssignment)
		}
	}

	// Convert the final assignments list to a set
	if len(resultAssignments) > 0 {
		tflog.Debug(ctx, "Setting final assignments", map[string]interface{}{
			"count": len(resultAssignments),
		})

		// Define attribute types for role assignment objects
		attrTypes := map[string]attr.Type{
			"id":                          types.StringType,
			"display_name":                types.StringType,
			"description":                 types.StringType,
			"admin_group_users_group_ids": types.SetType{ElemType: types.StringType},
			"scope_type":                  types.StringType,
			"resource_scopes":             types.SetType{ElemType: types.StringType},
		}

		assignmentsSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrTypes}, resultAssignments)
		if diags.HasError() {
			tflog.Error(ctx, "Error converting assignments to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})
		} else {
			roleDefinition.Assignments = assignmentsSet
		}
	} else if len(resultAssignments) == 0 && len(plannedAssignments) == 0 {
		// If no assignments found and none planned, set empty set
		attrTypes := map[string]attr.Type{
			"id":                          types.StringType,
			"display_name":                types.StringType,
			"description":                 types.StringType,
			"admin_group_users_group_ids": types.SetType{ElemType: types.StringType},
			"scope_type":                  types.StringType,
			"resource_scopes":             types.SetType{ElemType: types.StringType},
		}

		emptySet, _ := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrTypes}, []sharedmodels.RoleAssignmentResourceModel{})
		roleDefinition.Assignments = emptySet
	}

	tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
}
