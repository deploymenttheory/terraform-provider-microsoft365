package graphBetaRoleDefinition

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps a remote assignment to the Terraform state model
func MapRemoteAssignmentStateToTerraform(ctx context.Context, data *sharedmodels.RoleAssignmentResourceModel, remoteAssignments graphmodels.RoleAssignmentCollectionResponseable) {
	if remoteAssignments == nil {
		tflog.Debug(ctx, "Remote assignments response is nil")
		return
	}

	assignments := remoteAssignments.GetValue()
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found")
		return
	}

	// Take the first assignment
	assignment := assignments[0]
	if assignment == nil {
		tflog.Debug(ctx, "First assignment is nil")
		return
	}

	assignmentID := state.StringPtrToString(assignment.GetId())
	tflog.Debug(ctx, "Mapping remote assignment state to Terraform", map[string]interface{}{
		"assignmentId": assignmentID,
	})

	// Set basic properties
	data.ID = types.StringValue(assignmentID)
	data.DisplayName = types.StringValue(state.StringPtrToString(assignment.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(assignment.GetDescription()))

	// Convert scope members to set
	scopeMembers := assignment.GetScopeMembers()
	if len(scopeMembers) > 0 {
		scopeMembersSet, diags := types.SetValueFrom(ctx, types.StringType, scopeMembers)
		if !diags.HasError() {
			data.ScopeMembers = scopeMembersSet
		} else {
			tflog.Error(ctx, "Error converting scope members to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})
			emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
			data.ScopeMembers = emptySet
		}
	} else {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.ScopeMembers = emptySet
	}

	// Convert resource scopes to set
	resourceScopes := assignment.GetResourceScopes()
	if len(resourceScopes) > 0 {
		resourceScopesSet, diags := types.SetValueFrom(ctx, types.StringType, resourceScopes)
		if !diags.HasError() {
			data.ResourceScopes = resourceScopesSet
		} else {
			tflog.Error(ctx, "Error converting resource scopes to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})
			emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
			data.ResourceScopes = emptySet
		}
	} else {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.ResourceScopes = emptySet
	}

	// Handle scope type
	if scopeType := assignment.GetScopeType(); scopeType != nil {
		data.ScopeType = types.StringValue(scopeType.String())
	} else {
		data.ScopeType = types.StringNull()
	}

	tflog.Debug(ctx, "Finished mapping remote assignment state", map[string]interface{}{
		"assignmentId": assignmentID,
	})
}
