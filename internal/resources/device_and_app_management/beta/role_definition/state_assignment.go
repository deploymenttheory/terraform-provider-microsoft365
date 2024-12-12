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

	tflog.Debug(ctx, "Starting to map remote assignment state to Terraform state", map[string]interface{}{
		"assignmentId": state.StringPtrToString(assignment.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(assignment.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(assignment.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(assignment.GetDescription()))

	// Convert scope members
	if scopeMembers := assignment.GetScopeMembers(); len(scopeMembers) > 0 {
		data.ScopeMembers = make([]types.String, len(scopeMembers))
		for i, member := range scopeMembers {
			data.ScopeMembers[i] = types.StringValue(member)
		}
	} else {
		data.ScopeMembers = []types.String{}
	}

	// Convert resource scopes
	if resourceScopes := assignment.GetResourceScopes(); len(resourceScopes) > 0 {
		data.ResourceScopes = make([]types.String, len(resourceScopes))
		for i, scope := range resourceScopes {
			data.ResourceScopes[i] = types.StringValue(scope)
		}
	} else {
		data.ResourceScopes = []types.String{}
	}

	// Handle scope type
	if scopeType := assignment.GetScopeType(); scopeType != nil {
		data.ScopeType = types.StringValue(scopeType.String())
	} else {
		data.ScopeType = types.StringNull()
	}

	tflog.Debug(ctx, "Finished mapping remote assignment state to Terraform state", map[string]interface{}{
		"assignmentId": data.ID.ValueString(),
	})
}
