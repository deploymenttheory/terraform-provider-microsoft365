package graphBetaRoleDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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

	resourceID := state.StringPtrToString(remoteResource.GetId())
	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{
		"resourceId": resourceID,
	})

	// Set basic properties
	data.ID = types.StringValue(resourceID)
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.IsBuiltIn = state.BoolPtrToTypeBool(remoteResource.GetIsBuiltIn())
	data.IsBuiltInRoleDefinition = state.BoolPtrToTypeBool(remoteResource.GetIsBuiltInRoleDefinition())

	// Process rolePermissions
	rolePermissions := remoteResource.GetRolePermissions()
	if len(rolePermissions) > 0 {
		mappedPermissions := make([]RolePermissionResourceModel, 0, len(rolePermissions))

		for _, rp := range rolePermissions {
			permModel := RolePermissionResourceModel{}

			// Collect all allowed resource actions
			var allAllowedActions []string

			// Add actions from resourceActions field
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
	} else {
		data.RolePermissions = []RolePermissionResourceModel{}
	}

	// Map role scope tag IDs
	scopeTagIds := remoteResource.GetRoleScopeTagIds()
	if len(scopeTagIds) > 0 {
		scopeTagSet, diags := types.SetValueFrom(ctx, types.StringType, scopeTagIds)
		if !diags.HasError() {
			data.RoleScopeTagIds = scopeTagSet
		} else {
			tflog.Error(ctx, "Error converting scope tags to set", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})
			// Create empty set with StringType
			emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
			data.RoleScopeTagIds = emptySet
		}
	} else {
		// Create empty set with StringType
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.RoleScopeTagIds = emptySet
	}

	tflog.Debug(ctx, "Finished mapping remote state", map[string]interface{}{
		"resourceId": resourceID,
	})
}
