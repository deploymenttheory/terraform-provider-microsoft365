package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a RoleDefinition resource using data from the Terraform model.
func constructResource(ctx context.Context, data *RoleDefinitionResourceModel) (graphmodels.RoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewRoleDefinition()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.IsBuiltIn.IsNull() && !data.IsBuiltIn.IsUnknown() {
		isBuiltIn := data.IsBuiltIn.ValueBool()
		requestBody.SetIsBuiltIn(&isBuiltIn)
	}

	if !data.IsBuiltInRoleDefinition.IsNull() && !data.IsBuiltInRoleDefinition.IsUnknown() {
		isBuiltInRoleDefinition := data.IsBuiltInRoleDefinition.ValueBool()
		requestBody.SetIsBuiltInRoleDefinition(&isBuiltInRoleDefinition)
	}

	// Handle Permissions
	if len(data.Permissions) > 0 {
		permissions := constructRolePermissions(data.Permissions)
		requestBody.SetPermissions(permissions)
	}

	// Handle RolePermissions (same structure as Permissions)
	if len(data.RolePermissions) > 0 {
		rolePermissions := constructRolePermissions(data.RolePermissions)
		requestBody.SetRolePermissions(rolePermissions)
	}

	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		var roleScopeTagIds []string
		for _, id := range data.RoleScopeTagIds.Elements() {
			if idStr, ok := id.(types.String); ok {
				if !idStr.IsNull() && !idStr.IsUnknown() {
					roleScopeTagIds = append(roleScopeTagIds, idStr.ValueString())
				}
			}
		}
		requestBody.SetRoleScopeTagIds(roleScopeTagIds)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRolePermissions helper function to construct role permissions
func constructRolePermissions(permissions []RolePermissionResourceModel) []graphmodels.RolePermissionable {
	rolePermissions := make([]graphmodels.RolePermissionable, 0, len(permissions))

	for _, v := range permissions {
		rolePermission := graphmodels.NewRolePermission()

		if !v.Actions.IsNull() && !v.Actions.IsUnknown() {
			var actions []string
			for _, a := range v.Actions.Elements() {
				if actionStr, ok := a.(types.String); ok {
					if !actionStr.IsNull() && !actionStr.IsUnknown() {
						actions = append(actions, actionStr.ValueString())
					}
				}
			}
			rolePermission.SetActions(actions)
		}

		if len(v.ResourceActions) > 0 {
			resourceActions := make([]graphmodels.ResourceActionable, 0, len(v.ResourceActions))
			for _, ra := range v.ResourceActions {
				resourceAction := graphmodels.NewResourceAction()

				if !ra.AllowedResourceActions.IsNull() && !ra.AllowedResourceActions.IsUnknown() {
					var allowedActions []string
					for _, a := range ra.AllowedResourceActions.Elements() {
						if actionStr, ok := a.(types.String); ok {
							if !actionStr.IsNull() && !actionStr.IsUnknown() {
								allowedActions = append(allowedActions, actionStr.ValueString())
							}
						}
					}
					resourceAction.SetAllowedResourceActions(allowedActions)
				}

				if !ra.NotAllowedResourceActions.IsNull() && !ra.NotAllowedResourceActions.IsUnknown() {
					var notAllowedActions []string
					for _, a := range ra.NotAllowedResourceActions.Elements() {
						if actionStr, ok := a.(types.String); ok {
							if !actionStr.IsNull() && !actionStr.IsUnknown() {
								notAllowedActions = append(notAllowedActions, actionStr.ValueString())
							}
						}
					}
					resourceAction.SetNotAllowedResourceActions(notAllowedActions)
				}

				resourceActions = append(resourceActions, resourceAction)
			}
			rolePermission.SetResourceActions(resourceActions)
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	return rolePermissions
}
