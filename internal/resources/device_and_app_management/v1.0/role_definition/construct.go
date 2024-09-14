package graphroledefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func constructResource(ctx context.Context, data *RoleDefinitionResourceModel) (models.RoleDefinitionable, error) {
	tflog.Debug(ctx, "Constructing RoleDefinition resource")
	construct.DebugPrintStruct(ctx, "Constructed RoleDefinition Resource from model", data)

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

	return roleDef, nil
}
