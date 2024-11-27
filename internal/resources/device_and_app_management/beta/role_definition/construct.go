package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
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

	if len(data.RolePermissions) > 0 {
		rolePermissions := make([]graphmodels.RolePermissionable, 0, len(data.RolePermissions))
		for _, v := range data.RolePermissions {
			rolePermission := graphmodels.NewRolePermission()

			if len(v.Actions) > 0 {
				actions := make([]string, 0, len(v.Actions))
				for _, a := range v.Actions {
					if !a.IsNull() && !a.IsUnknown() {
						actions = append(actions, a.ValueString())
					}
				}
				rolePermission.SetActions(actions)
			}

			if len(v.ResourceActions) > 0 {
				resourceActions := make([]graphmodels.ResourceActionable, 0, len(v.ResourceActions))
				for _, ra := range v.ResourceActions {
					resourceAction := graphmodels.NewResourceAction()

					var allowedActions []string
					for _, a := range ra.AllowedResourceActions {
						if !a.IsNull() && !a.IsUnknown() {
							allowedActions = append(allowedActions, a.ValueString())
						}
					}
					resourceAction.SetAllowedResourceActions(allowedActions)

					var notAllowedActions []string
					for _, a := range ra.NotAllowedResourceActions {
						if !a.IsNull() && !a.IsUnknown() {
							notAllowedActions = append(notAllowedActions, a.ValueString())
						}
					}
					resourceAction.SetNotAllowedResourceActions(notAllowedActions)

					resourceActions = append(resourceActions, resourceAction)
				}
				rolePermission.SetResourceActions(resourceActions)
			}

			rolePermissions = append(rolePermissions, rolePermission)
		}
		requestBody.SetRolePermissions(rolePermissions)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, id := range data.RoleScopeTagIds {
			if !id.IsNull() && !id.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, id.ValueString())
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
