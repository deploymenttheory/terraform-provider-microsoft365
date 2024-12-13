package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructAssignment(ctx context.Context, data *RoleDefinitionResourceModel) (graphmodels.RoleAssignmentable, error) {
	tflog.Debug(ctx, "Constructing role assignment")

	requestBody := graphmodels.NewRoleAssignment()

	if !data.Assignments.DisplayName.IsNull() {
		displayName := data.Assignments.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Assignments.Description.IsNull() {
		description := data.Assignments.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if len(data.Assignments.ScopeMembers) > 0 {
		var scopeMembers []string
		for _, member := range data.Assignments.ScopeMembers {
			if !member.IsNull() && !member.IsUnknown() {
				scopeMembers = append(scopeMembers, member.ValueString())
			}
		}
		requestBody.SetScopeMembers(scopeMembers)
	}

	if len(data.Assignments.ResourceScopes) > 0 {
		var resourceScopes []string
		for _, scope := range data.Assignments.ResourceScopes {
			if !scope.IsNull() && !scope.IsUnknown() {
				resourceScopes = append(resourceScopes, scope.ValueString())
			}
		}
		requestBody.SetResourceScopes(resourceScopes)
	}

	if !data.Assignments.ScopeType.IsNull() && !data.Assignments.ScopeType.IsUnknown() {
		scopeTypeStr := data.Assignments.ScopeType.ValueString()
		scopeTypeVal, err := graphmodels.ParseRoleAssignmentScopeType(scopeTypeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing role assignment scope type: %v", err)
		}
		if scopeTypeVal != nil {
			scopeType := scopeTypeVal.(*graphmodels.RoleAssignmentScopeType)
			requestBody.SetScopeType(scopeType)
		}
	}
	if err := construct.DebugLogGraphObject(ctx, "Role Assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing role assignment")
	return requestBody, nil
}
