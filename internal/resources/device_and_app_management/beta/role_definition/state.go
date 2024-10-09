package graphbetaroledefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *RoleDefinitionResourceModel, remoteResource graphmodels.RoleDefinitionable) {
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

	// Handle RolePermissions
	rolePermissions := remoteResource.GetRolePermissions()
	if len(rolePermissions) > 0 {
		data.RolePermissions = make([]RolePermissionResourceModel, len(rolePermissions))
		for i, rp := range rolePermissions {
			resourceActions := rp.GetResourceActions()
			if len(resourceActions) > 0 {
				data.RolePermissions[i].ResourceActions = make([]ResourceActionResourceModel, len(resourceActions))
				for j, ra := range resourceActions {
					data.RolePermissions[i].ResourceActions[j] = ResourceActionResourceModel{
						AllowedResourceActions:    state.SliceToTypeStringSlice(ra.GetAllowedResourceActions()),
						NotAllowedResourceActions: state.SliceToTypeStringSlice(ra.GetNotAllowedResourceActions()),
					}
				}
			}
		}
	} else {
		data.RolePermissions = []RolePermissionResourceModel{}
	}

	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
