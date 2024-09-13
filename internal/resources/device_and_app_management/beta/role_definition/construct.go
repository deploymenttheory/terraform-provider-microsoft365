package graphbetaroledefinition

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *RoleDefinitionResourceModel) (models.RoleDefinitionable, error) {
	tflog.Debug(ctx, "Constructing RoleDefinition resource")

	roleDef := models.NewRoleDefinition()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		roleDef.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		roleDef.SetDescription(&description)
	}

	if !data.IsBuiltIn.IsNull() && !data.IsBuiltIn.IsUnknown() {
		isBuiltIn := data.IsBuiltIn.ValueBool()
		roleDef.SetIsBuiltIn(&isBuiltIn)
	}

	if !data.IsBuiltInRoleDefinition.IsNull() && !data.IsBuiltInRoleDefinition.IsUnknown() {
		isBuiltInRoleDefinition := data.IsBuiltInRoleDefinition.ValueBool()
		roleDef.SetIsBuiltInRoleDefinition(&isBuiltInRoleDefinition)
	}

	if len(data.RolePermissions) > 0 {
		rolePermissions := make([]models.RolePermissionable, 0, len(data.RolePermissions))
		for _, v := range data.RolePermissions {
			rolePermission := models.NewRolePermission()

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
				resourceActions := make([]models.ResourceActionable, 0, len(v.ResourceActions))
				for _, ra := range v.ResourceActions {
					resourceAction := models.NewResourceAction()

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
		roleDef.SetRolePermissions(rolePermissions)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, id := range data.RoleScopeTagIds {
			if !id.IsNull() && !id.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, id.ValueString())
			}
		}
		roleDef.SetRoleScopeTagIds(roleScopeTagIds)
	}

	debugPrintRoleDefinition(ctx, roleDef)

	return roleDef, nil
}

func debugPrintRoleDefinition(ctx context.Context, roleDef models.RoleDefinitionable) {
	roleDefMap := map[string]interface{}{
		"id":                          roleDef.GetId(),
		"display_name":                roleDef.GetDisplayName(),
		"description":                 roleDef.GetDescription(),
		"is_built_in":                 roleDef.GetIsBuiltIn(),
		"is_built_in_role_definition": roleDef.GetIsBuiltInRoleDefinition(),
		"role_scope_tag_ids":          roleDef.GetRoleScopeTagIds(),
		"role_permissions":            debugMapRolePermissions(roleDef.GetRolePermissions()),
	}

	roleDefJSON, err := json.MarshalIndent(roleDefMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling RoleDefinition to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed RoleDefinition resource", map[string]interface{}{
		"roleDefinition": string(roleDefJSON),
	})
}

func debugMapRolePermissions(permissions []models.RolePermissionable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(permissions))
	for i, perm := range permissions {
		result[i] = map[string]interface{}{
			"actions":          perm.GetActions(),
			"resource_actions": debugMapResourceActions(perm.GetResourceActions()),
		}
	}
	return result
}

func debugMapResourceActions(actions []models.ResourceActionable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(actions))
	for i, action := range actions {
		result[i] = map[string]interface{}{
			"allowed_resource_actions":     action.GetAllowedResourceActions(),
			"not_allowed_resource_actions": action.GetNotAllowedResourceActions(),
		}
	}
	return result
}
