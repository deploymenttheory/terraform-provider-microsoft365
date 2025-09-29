package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote RoleDefinition state to Terraform
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleDefinitionResourceModel, remoteResource graphmodels.RoleDefinitionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	resourceID := convert.GraphToFrameworkString(remoteResource.GetId()).ValueString()

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]any{
		"resourceId": resourceID,
	})

	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsBuiltIn = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())
	data.IsBuiltInRoleDefinition = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltInRoleDefinition())

	rolePermissions := remoteResource.GetRolePermissions()

	tflog.Debug(ctx, fmt.Sprintf("API returned %d rolePermissions", len(rolePermissions)))

	if len(rolePermissions) > 0 {
		mappedPermissions := make([]RolePermissionResourceModel, 0, len(rolePermissions))

		for _, rp := range rolePermissions {
			permModel := RolePermissionResourceModel{}

			var allAllowedActions []string

			resourceActions := rp.GetResourceActions()
			tflog.Debug(ctx, fmt.Sprintf("Role permission has %d resourceActions", len(resourceActions)))

			for _, ra := range resourceActions {
				allowedActions := ra.GetAllowedResourceActions()
				if len(allowedActions) > 0 {
					tflog.Debug(ctx, fmt.Sprintf("Found %d allowed resource actions: %v", len(allowedActions), allowedActions))
					allAllowedActions = append(allAllowedActions, allowedActions...)
				}
			}

			// Add actions from the actions field (if they exist and aren't already in the list)
			apiActions := rp.GetActions()
			if len(apiActions) > 0 {
				tflog.Debug(ctx, fmt.Sprintf("Found %d actions: %v", len(apiActions), apiActions))
				for _, action := range apiActions {
					// Check if action is already in allAllowedActions
					found := false
					for _, existing := range allAllowedActions {
						if existing == action {
							found = true
							break
						}
					}
					if !found {
						allAllowedActions = append(allAllowedActions, action)
					}
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Total actions collected: %d", len(allAllowedActions)))

			// Create set with proper element type
			if len(allAllowedActions) > 0 {
				allowedActionsSet, diags := types.SetValueFrom(ctx, types.StringType, allAllowedActions)
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

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
