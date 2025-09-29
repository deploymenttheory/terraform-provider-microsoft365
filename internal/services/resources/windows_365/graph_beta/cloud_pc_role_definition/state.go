package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote UnifiedRoleDefinition state to Terraform
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleDefinitionResourceModel, remoteResource graphmodels.UnifiedRoleDefinitionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	resourceID := data.ID.ValueString()

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]any{
		"resourceId": resourceID,
	})

	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsBuiltIn = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())
	data.IsBuiltInRoleDefinition = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())

	rolePermissions := remoteResource.GetRolePermissions()

	tflog.Debug(ctx, fmt.Sprintf("API returned %d rolePermissions", len(rolePermissions)))

	if len(rolePermissions) > 0 {
		mappedPermissions := make([]RolePermissionResourceModel, 0, len(rolePermissions))

		for _, rp := range rolePermissions {
			permModel := RolePermissionResourceModel{}

			// For UnifiedRolePermissionable, get allowed resource actions directly
			allowedActions := rp.GetAllowedResourceActions()
			tflog.Debug(ctx, fmt.Sprintf("Found %d allowed resource actions: %v", len(allowedActions), allowedActions))

			tflog.Debug(ctx, fmt.Sprintf("Total actions collected: %d", len(allowedActions)))

			// Create set with proper element type
			if len(allowedActions) > 0 {
				allowedActionsSet, diags := types.SetValueFrom(ctx, types.StringType, allowedActions)
				if !diags.HasError() {
					permModel.AllowedResourceActions = allowedActionsSet
				} else {
					tflog.Error(ctx, "Error converting allowed resource actions to set", map[string]any{
						"error": diags.Errors()[0].Detail(),
					})
					// Create empty set with StringType
					permModel.AllowedResourceActions, _ = types.SetValueFrom(ctx, types.StringType, []string{})
				}
			} else {
				// Create empty set with StringType
				emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
				permModel.AllowedResourceActions = emptySet
			}

			mappedPermissions = append(mappedPermissions, permModel)
		}

		data.RolePermissions = mappedPermissions
		tflog.Debug(ctx, fmt.Sprintf("Set %d role permissions in state", len(mappedPermissions)))
	} else {
		tflog.Warn(ctx, "No role permissions returned from API - this may indicate the permissions are stored elsewhere or the API call is incomplete")
	}
	// Note: If no role permissions are returned from API, we don't set data.RolePermissions at all,
	// leaving it as whatever was in the original state/plan to maintain consistency

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
