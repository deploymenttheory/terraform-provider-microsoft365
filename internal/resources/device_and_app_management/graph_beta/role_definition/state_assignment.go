package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"sort"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state" // Import the state helpers
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// using the API response and state helpers, ensuring stable ordering.
func MapRemoteAssignmentStateToTerraform(ctx context.Context, roleDefinition *RoleDefinitionResourceModel, remoteAssignments graphmodels.RoleAssignmentCollectionResponseable) {
	tflog.Debug(ctx, "Starting MapRemoteAssignmentStateToTerraform")

	// Define the object type for assignments (remains the same)
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

	// If no assignments remotely, set to null and return (remains the same)
	if remoteAssignments == nil || remoteAssignments.GetValue() == nil || len(remoteAssignments.GetValue()) == 0 {
		tflog.Debug(ctx, "No remote assignments to process, setting assignments to null")
		roleDefinition.Assignments = types.SetNull(assignmentObjectType)
		tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
		return
	}

	// Process remote assignments
	assignList := remoteAssignments.GetValue()

	// Build objects for set value
	var objects []attr.Value
	for _, assignment := range assignList {
		if assignment == nil {
			continue
		}

		// Use helpers for string pointer conversions
		id := state.StringPointerValue(assignment.GetId())             // Use helper
		dName := state.StringPointerValue(assignment.GetDisplayName()) // Use helper
		desc := state.StringPointerValue(assignment.GetDescription())  // Use helper

		// Sort members and scopes for consistent ordering (keep sorting)
		var sortedMembers, sortedScopes []string
		if m := assignment.GetScopeMembers(); len(m) > 0 {
			sortedMembers = append(sortedMembers, m...)
			sort.Strings(sortedMembers)
		}
		if s := assignment.GetResourceScopes(); len(s) > 0 {
			sortedScopes = append(sortedScopes, s...)
			sort.Strings(sortedScopes)
		}

		// Determine scope type, defaulting to resourceScope
		var scopeTypeVal types.String
		if st := assignment.GetScopeType(); st != nil {
			scopeTypeVal = state.StringValue(st.String()) // Use helper
		} else {
			scopeTypeVal = state.StringValue("resourceScope") // Use helper for consistency
		}

		// Use helper to create sets from sorted slices
		// The helper returns SetNull on error or if the input slice is empty/nil.
		membersSet := state.StringSliceToSet(ctx, sortedMembers)
		scopesSet := state.StringSliceToSet(ctx, sortedScopes)

		// Create the object attributes map using the results from helpers
		objAttrs := map[string]attr.Value{
			"id":                          id,
			"display_name":                dName,
			"description":                 desc,
			"admin_group_users_group_ids": membersSet, // Directly use the result from StringSliceToSet
			"scope_type":                  scopeTypeVal,
			"resource_scopes":             scopesSet, // Directly use the result from StringSliceToSet
		}

		// Create the object (error handling remains similar)
		obj, diags := types.ObjectValue(assignmentObjectType.AttrTypes, objAttrs)
		if diags.HasError() {
			for _, d := range diags.Errors() {
				tflog.Error(ctx, "Error creating assignment object", map[string]interface{}{
					"error": d.Detail(),
					"id":    id.ValueString(), // Use ValueString() as id is now types.String
				})
			}
			continue // Skip this assignment if object creation fails
		}

		objects = append(objects, obj)
	}

	// Convert object slice to set (error handling remains similar)
	setVal, diags := types.SetValue(assignmentObjectType, objects)
	if diags.HasError() {
		tflog.Error(ctx, "Error converting assignments to set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
		// Set to null on error
		roleDefinition.Assignments = types.SetNull(assignmentObjectType)
	} else {
		roleDefinition.Assignments = setVal
		tflog.Debug(ctx, fmt.Sprintf("Mapped %d assignments into state", len(objects)))
	}

	tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
}
