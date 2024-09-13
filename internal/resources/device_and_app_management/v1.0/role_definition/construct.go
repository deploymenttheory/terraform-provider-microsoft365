package graphroledefinition

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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

	if len(data.RolePermissions) > 0 {
		rolePermissions := make([]models.RolePermissionable, 0, len(data.RolePermissions))
		for _, v := range data.RolePermissions {
			rolePermission := models.NewRolePermission()

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

	// Debug logging
	debugPrintRequestBody(ctx, roleDef)

	return roleDef, nil
}

func debugPrintRequestBody(ctx context.Context, roleDef models.RoleDefinitionable) {
	requestMap := map[string]interface{}{
		"displayName":     roleDef.GetDisplayName(),
		"description":     roleDef.GetDescription(),
		"isBuiltIn":       roleDef.GetIsBuiltIn(),
		"rolePermissions": debugMapRolePermissions(roleDef.GetRolePermissions()),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed RoleDefinition resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}

func debugMapRolePermissions(permissions []models.RolePermissionable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(permissions))
	for i, perm := range permissions {
		result[i] = map[string]interface{}{
			"resourceActions": debugMapResourceActions(perm.GetResourceActions()),
		}
	}
	return result
}

func debugMapResourceActions(actions []models.ResourceActionable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(actions))
	for i, action := range actions {
		result[i] = map[string]interface{}{
			"allowedResourceActions":    action.GetAllowedResourceActions(),
			"notAllowedResourceActions": action.GetNotAllowedResourceActions(),
		}
	}
	return result
}
