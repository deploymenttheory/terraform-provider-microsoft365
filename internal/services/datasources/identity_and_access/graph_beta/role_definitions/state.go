package graphBetaRoleDefinitions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToDataSource(ctx context.Context, roleDefinition graphmodels.UnifiedRoleDefinitionable) RoleDefinitionModel {
	model := RoleDefinitionModel{}

	if roleDefinition.GetId() != nil {
		model.ID = types.StringValue(*roleDefinition.GetId())
	} else {
		model.ID = types.StringNull()
	}

	if roleDefinition.GetDescription() != nil {
		model.Description = types.StringValue(*roleDefinition.GetDescription())
	} else {
		model.Description = types.StringNull()
	}

	if roleDefinition.GetDisplayName() != nil {
		model.DisplayName = types.StringValue(*roleDefinition.GetDisplayName())
	} else {
		model.DisplayName = types.StringNull()
	}

	if roleDefinition.GetIsBuiltIn() != nil {
		model.IsBuiltIn = types.BoolValue(*roleDefinition.GetIsBuiltIn())
	} else {
		model.IsBuiltIn = types.BoolNull()
	}

	if roleDefinition.GetIsEnabled() != nil {
		model.IsEnabled = types.BoolValue(*roleDefinition.GetIsEnabled())
	} else {
		model.IsEnabled = types.BoolNull()
	}

	if roleDefinition.GetIsPrivileged() != nil {
		model.IsPrivileged = types.BoolValue(*roleDefinition.GetIsPrivileged())
	} else {
		model.IsPrivileged = types.BoolNull()
	}

	if resourceScopes := roleDefinition.GetResourceScopes(); resourceScopes != nil {
		model.ResourceScopes = make([]types.String, len(resourceScopes))
		for i, scope := range resourceScopes {
			model.ResourceScopes[i] = types.StringValue(scope)
		}
	} else {
		model.ResourceScopes = []types.String{}
	}

	if roleDefinition.GetTemplateId() != nil {
		model.TemplateID = types.StringValue(*roleDefinition.GetTemplateId())
	} else {
		model.TemplateID = types.StringNull()
	}

	if roleDefinition.GetVersion() != nil {
		model.Version = types.StringValue(*roleDefinition.GetVersion())
	} else {
		model.Version = types.StringNull()
	}

	if rolePermissions := roleDefinition.GetRolePermissions(); rolePermissions != nil {
		model.RolePermissions = make([]RolePermissionModel, len(rolePermissions))
		for i, perm := range rolePermissions {
			permModel := RolePermissionModel{}

			if allowedActions := perm.GetAllowedResourceActions(); allowedActions != nil {
				permModel.AllowedResourceActions = make([]types.String, len(allowedActions))
				for j, action := range allowedActions {
					permModel.AllowedResourceActions[j] = types.StringValue(action)
				}
			} else {
				permModel.AllowedResourceActions = []types.String{}
			}

			if perm.GetCondition() != nil {
				permModel.Condition = types.StringValue(*perm.GetCondition())
			} else {
				permModel.Condition = types.StringNull()
			}

			model.RolePermissions[i] = permModel
		}
	} else {
		model.RolePermissions = []RolePermissionModel{}
	}

	tflog.Debug(ctx, "Successfully mapped role definition to data source model", map[string]any{
		"id":          model.ID.ValueString(),
		"displayName": model.DisplayName.ValueString(),
	})

	return model
}
