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

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{
		"resourceId": resourceID,
	})

	// Set basic properties
	data.ID = types.StringValue(resourceID)
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsBuiltIn = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())
	data.IsBuiltInRoleDefinition = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltInRoleDefinition())

	rolePermissions := remoteResource.GetRolePermissions()
	if len(rolePermissions) > 0 {
		mappedPermissions := make([]RolePermissionResourceModel, 0, len(rolePermissions))

		for _, rp := range rolePermissions {
			permModel := RolePermissionResourceModel{}

			var allAllowedActions []string

			resourceActions := rp.GetResourceActions()
			for _, ra := range resourceActions {
				allowedActions := ra.GetAllowedResourceActions()
				if len(allowedActions) > 0 {
					allAllowedActions = append(allAllowedActions, allowedActions...)
				}
			}

			// Add actions from the actions field (if they exist and aren't already in the list)
			apiActions := rp.GetActions()
			if len(apiActions) > 0 {
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

			// Create set with proper element type
			if len(allAllowedActions) > 0 {
				allowedActionsSet, diags := types.SetValueFrom(ctx, types.StringType, allAllowedActions)
				if !diags.HasError() {
					permModel.AllowedResourceActions = allowedActionsSet
				} else {
					tflog.Error(ctx, "Error converting allowed resource actions to set", map[string]interface{}{
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
	}
	// Note: If no role permissions are returned from API, we don't set data.RolePermissions at all,
	// leaving it as whatever was in the original state/plan to maintain consistency

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
