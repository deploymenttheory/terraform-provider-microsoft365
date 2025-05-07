package graphBetaRoleDefinitionAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps a remote role assignment to the Terraform resource model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleDefinitionAssignmentResourceModel, assignment graphmodels.DeviceAndAppManagementRoleAssignmentable) {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return
	}

	assignmentID := state.StringPtrToString(assignment.GetId())
	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{
		"assignmentId": assignmentID,
	})

	data.ID = types.StringValue(assignmentID)
	data.DisplayName = types.StringValue(state.StringPtrToString(assignment.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(assignment.GetDescription()))
	data.ScopeType = state.EnumPtrToTypeString(assignment.GetScopeType())

	if members := assignment.GetScopeMembers(); len(members) > 0 {
		scopeMembers, diags := types.SetValueFrom(ctx, types.StringType, members)
		if !diags.HasError() {
			data.ScopeMembers = scopeMembers
		} else {
			tflog.Error(ctx, "Error converting scope members to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
				"id":    assignmentID,
			})
		}
	}

	// Set resource scopes
	if scopes := assignment.GetResourceScopes(); len(scopes) > 0 {
		resourceScopes, diags := types.SetValueFrom(ctx, types.StringType, scopes)
		if !diags.HasError() {
			data.ResourceScopes = resourceScopes
		} else {
			tflog.Error(ctx, "Error converting resource scopes to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
				"id":    assignmentID,
			})
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state", map[string]interface{}{
		"assignmentId": assignmentID,
	})
}
