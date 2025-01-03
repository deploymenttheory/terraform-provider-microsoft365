package graphBetaRoleDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a RoleDefinitionResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleDefinitionResourceModel, remoteResource graphmodels.RoleDefinitionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.IsBuiltIn = state.BoolPtrToTypeBool(remoteResource.GetIsBuiltIn())
	data.IsBuiltInRoleDefinition = state.BoolPtrToTypeBool(remoteResource.GetIsBuiltInRoleDefinition())

	permissions := remoteResource.GetPermissions()
	if len(permissions) > 0 {
		data.Permissions = make([]RolePermissionResourceModel, len(permissions))
		for i, p := range permissions {
			actionSlice := state.SliceToTypeStringSlice(p.GetActions())
			actionsSet, diags := types.SetValueFrom(ctx, types.StringType, actionSlice)
			if diags.HasError() {
				tflog.Error(ctx, "Error converting actions to set", map[string]interface{}{
					"error": diags.Errors()[0].Detail(),
				})
				continue
			}

			resourceActions := p.GetResourceActions()
			raModels := make([]ResourceActionResourceModel, len(resourceActions))
			for j, ra := range resourceActions {
				allowedSlice := state.SliceToTypeStringSlice(ra.GetAllowedResourceActions())
				allowedSet, diags := types.SetValueFrom(ctx, types.StringType, allowedSlice)
				if diags.HasError() {
					continue
				}

				notAllowedSlice := state.SliceToTypeStringSlice(ra.GetNotAllowedResourceActions())
				notAllowedSet, diags := types.SetValueFrom(ctx, types.StringType, notAllowedSlice)
				if diags.HasError() {
					continue
				}

				raModels[j] = ResourceActionResourceModel{
					AllowedResourceActions:    allowedSet,
					NotAllowedResourceActions: notAllowedSet,
				}
			}

			data.Permissions[i] = RolePermissionResourceModel{
				Actions:         actionsSet,
				ResourceActions: raModels,
			}
		}
	} else {
		data.Permissions = []RolePermissionResourceModel{}
	}

	rolePermissions := remoteResource.GetRolePermissions()
	if len(rolePermissions) > 0 {
		data.RolePermissions = make([]RolePermissionResourceModel, len(rolePermissions))
		for i, rp := range rolePermissions {
			actionSlice := state.SliceToTypeStringSlice(rp.GetActions())
			actionsSet, diags := types.SetValueFrom(ctx, types.StringType, actionSlice)
			if diags.HasError() {
				continue
			}

			resourceActions := rp.GetResourceActions()
			raModels := make([]ResourceActionResourceModel, len(resourceActions))
			for j, ra := range resourceActions {
				allowedSlice := state.SliceToTypeStringSlice(ra.GetAllowedResourceActions())
				allowedSet, diags := types.SetValueFrom(ctx, types.StringType, allowedSlice)
				if diags.HasError() {
					continue
				}

				notAllowedSlice := state.SliceToTypeStringSlice(ra.GetNotAllowedResourceActions())
				notAllowedSet, diags := types.SetValueFrom(ctx, types.StringType, notAllowedSlice)
				if diags.HasError() {
					continue
				}

				raModels[j] = ResourceActionResourceModel{
					AllowedResourceActions:    allowedSet,
					NotAllowedResourceActions: notAllowedSet,
				}
			}

			data.RolePermissions[i] = RolePermissionResourceModel{
				Actions:         actionsSet,
				ResourceActions: raModels,
			}
		}
	} else {
		data.RolePermissions = []RolePermissionResourceModel{}
	}

	// Convert RoleScopeTagIds to Set
	scopeTagSlice := state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	scopeTagSet, diags := types.SetValueFrom(ctx, types.StringType, scopeTagSlice)
	if diags.HasError() {
		tflog.Error(ctx, "Error converting scope tags to set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
	} else {
		data.RoleScopeTagIds = scopeTagSet
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
