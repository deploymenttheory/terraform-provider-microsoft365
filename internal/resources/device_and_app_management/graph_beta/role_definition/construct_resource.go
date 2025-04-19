package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a RoleDefinition resource using data from the Terraform model.
// This implementation aligns with the Microsoft Graph example by consolidating all permissions
// into a single rolePermission with a single resourceAction.
func constructResource(ctx context.Context, data *RoleDefinitionResourceModel) (graphmodels.RoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewRoleDefinition()

	// Set basic properties
	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.IsBuiltIn.IsNull() && !data.IsBuiltIn.IsUnknown() {
		isBuiltIn := data.IsBuiltIn.ValueBool()
		requestBody.SetIsBuiltIn(&isBuiltIn)
	}

	if !data.IsBuiltInRoleDefinition.IsNull() && !data.IsBuiltInRoleDefinition.IsUnknown() {
		isBuiltInRoleDefinition := data.IsBuiltInRoleDefinition.ValueBool()
		requestBody.SetIsBuiltInRoleDefinition(&isBuiltInRoleDefinition)
	}

	// Create a single rolePermission with all allowed resource actions
	rolePermission := graphmodels.NewRolePermission()
	resourceAction := graphmodels.NewResourceAction()

	// Collect all allowed and not allowed resource actions
	var allowedResourceActions []string
	var notAllowedResourceActions []string

	// Process permissions from data.Permissions
	for _, perm := range data.Permissions {
		// Add actions from the actions set
		if !perm.Actions.IsNull() && !perm.Actions.IsUnknown() {
			for _, a := range perm.Actions.Elements() {
				if actionStr, ok := a.(types.String); ok {
					if !actionStr.IsNull() && !actionStr.IsUnknown() {
						allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
					}
				}
			}
		}

		// Add allowed/not allowed resource actions
		for _, ra := range perm.ResourceActions {
			if !ra.AllowedResourceActions.IsNull() && !ra.AllowedResourceActions.IsUnknown() {
				for _, a := range ra.AllowedResourceActions.Elements() {
					if actionStr, ok := a.(types.String); ok {
						if !actionStr.IsNull() && !actionStr.IsUnknown() {
							allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
						}
					}
				}
			}

			if !ra.NotAllowedResourceActions.IsNull() && !ra.NotAllowedResourceActions.IsUnknown() {
				for _, a := range ra.NotAllowedResourceActions.Elements() {
					if actionStr, ok := a.(types.String); ok {
						if !actionStr.IsNull() && !actionStr.IsUnknown() {
							notAllowedResourceActions = append(notAllowedResourceActions, actionStr.ValueString())
						}
					}
				}
			}
		}
	}

	// Process permissions from data.RolePermissions
	for _, perm := range data.RolePermissions {
		// Add actions from the actions set
		if !perm.Actions.IsNull() && !perm.Actions.IsUnknown() {
			for _, a := range perm.Actions.Elements() {
				if actionStr, ok := a.(types.String); ok {
					if !actionStr.IsNull() && !actionStr.IsUnknown() {
						allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
					}
				}
			}
		}

		// Add allowed/not allowed resource actions
		for _, ra := range perm.ResourceActions {
			if !ra.AllowedResourceActions.IsNull() && !ra.AllowedResourceActions.IsUnknown() {
				for _, a := range ra.AllowedResourceActions.Elements() {
					if actionStr, ok := a.(types.String); ok {
						if !actionStr.IsNull() && !actionStr.IsUnknown() {
							allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
						}
					}
				}
			}

			if !ra.NotAllowedResourceActions.IsNull() && !ra.NotAllowedResourceActions.IsUnknown() {
				for _, a := range ra.NotAllowedResourceActions.Elements() {
					if actionStr, ok := a.(types.String); ok {
						if !actionStr.IsNull() && !actionStr.IsUnknown() {
							notAllowedResourceActions = append(notAllowedResourceActions, actionStr.ValueString())
						}
					}
				}
			}
		}
	}

	// Set allowed resource actions
	resourceAction.SetAllowedResourceActions(allowedResourceActions)

	// Set not allowed resource actions if any
	if len(notAllowedResourceActions) > 0 {
		resourceAction.SetNotAllowedResourceActions(notAllowedResourceActions)
	}

	// Create resourceActions array with the single resourceAction
	resourceActions := []graphmodels.ResourceActionable{
		resourceAction,
	}
	rolePermission.SetResourceActions(resourceActions)

	// Create rolePermissions array with the single rolePermission
	rolePermissions := []graphmodels.RolePermissionable{
		rolePermission,
	}
	requestBody.SetRolePermissions(rolePermissions)

	// Set role scope tag IDs
	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		var roleScopeTagIds []string
		for _, id := range data.RoleScopeTagIds.Elements() {
			if idStr, ok := id.(types.String); ok {
				if !idStr.IsNull() && !idStr.IsUnknown() {
					roleScopeTagIds = append(roleScopeTagIds, idStr.ValueString())
				}
			}
		}
		requestBody.SetRoleScopeTagIds(roleScopeTagIds)
	}

	// Debug log the constructed object
	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
