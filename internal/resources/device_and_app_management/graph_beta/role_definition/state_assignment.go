package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"sort"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// using only the API response, ensuring stable ordering.
// MapRemoteAssignmentStateToTerraform maps remote assignments to the Terraform state model
// using only the API response, ensuring stable ordering.
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

	// If no assignments remotely, set to null and return
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

		// Extract fields
		id := state.StringPtrToString(assignment.GetId())
		dName := state.StringPtrToString(assignment.GetDisplayName())

		// Sort members and scopes for consistent ordering
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
		var scopeTypeVal string
		if st := assignment.GetScopeType(); st != nil {
			scopeTypeVal = st.String()
		} else {
			scopeTypeVal = "resourceScope"
		}

		// Create the members set
		membersSet, diags := types.SetValueFrom(ctx, types.StringType, sortedMembers)
		if diags.HasError() {
			tflog.Error(ctx, "Error converting scope members to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
				"id":    id,
			})
			continue
		}

		// Create the resource scopes set
		scopesSet, diags := types.SetValueFrom(ctx, types.StringType, sortedScopes)
		if diags.HasError() {
			tflog.Error(ctx, "Error converting resource scopes to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
				"id":    id,
			})
			continue
		}

		// Create the object attributes map
		objAttrs := map[string]attr.Value{
			"id":                          types.StringValue(id),
			"display_name":                types.StringValue(dName),
			"description":                 types.StringValue(state.StringPtrToString(assignment.GetDescription())),
			"admin_group_users_group_ids": membersSet,
			"scope_type":                  types.StringValue(scopeTypeVal),
			"resource_scopes":             scopesSet,
		}

		// Create the object
		obj, diags := types.ObjectValue(assignmentObjectType.AttrTypes, objAttrs)
		if diags.HasError() {
			for _, d := range diags.Errors() {
				tflog.Error(ctx, "Error creating assignment object", map[string]interface{}{
					"error": d.Detail(),
					"id":    id,
				})
			}
			continue
		}

		objects = append(objects, obj)
	}

	// Convert to set
	setVal, diags := types.SetValue(assignmentObjectType, objects)
	if diags.HasError() {
		tflog.Error(ctx, "Error converting assignments to set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
		roleDefinition.Assignments = types.SetNull(assignmentObjectType)
	} else {
		roleDefinition.Assignments = setVal
		tflog.Debug(ctx, fmt.Sprintf("Mapped %d assignments into state", len(objects)))
	}

	tflog.Debug(ctx, "Finished MapRemoteAssignmentStateToTerraform")
}
