package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a RoleDefinition resource using data from the Terraform model.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleDefinitionResourceModel, resp interface{}, readPermissions []string, isUpdate bool) (graphmodels.RoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewRoleDefinition()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	// For updates, don't set read-only or immutable properties
	if !isUpdate {
		if !data.IsBuiltIn.IsNull() {
			isBuiltIn := data.IsBuiltIn.ValueBool()
			requestBody.SetIsBuiltIn(&isBuiltIn)
		}

		if !data.IsBuiltInRoleDefinition.IsNull() {
			isBuiltInRoleDefinition := data.IsBuiltInRoleDefinition.ValueBool()
			requestBody.SetIsBuiltInRoleDefinition(&isBuiltInRoleDefinition)
		}
	}

	rolePermission := graphmodels.NewRolePermission()
	resourceAction := graphmodels.NewResourceAction()

	allowedResourceActions := []string{}

	for _, perm := range data.RolePermissions {
		if !perm.AllowedResourceActions.IsNull() {
			for _, a := range perm.AllowedResourceActions.Elements() {
				if actionStr, ok := a.(types.String); ok && !actionStr.IsNull() {
					allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
				}
			}
		}
	}

	validatedPermissions, err := validateRolePermissions(ctx, client, allowedResourceActions, resp, readPermissions)
	if err != nil {
		return nil, err
	}

	resourceAction.SetAllowedResourceActions(validatedPermissions)
	resourceActions := []graphmodels.ResourceActionable{resourceAction}
	rolePermission.SetResourceActions(resourceActions)
	rolePermissions := []graphmodels.RolePermissionable{rolePermission}
	requestBody.SetRolePermissions(rolePermissions)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

// validateRolePermissions validates that all provided role permissions exist in the list of available operations
func validateRolePermissions(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, permissions []string, resp interface{}, readPermissions []string) ([]string, error) {
	tflog.Debug(ctx, "Validating Intune role permissions against available Intune resource operations")

	// Get the list of available resource operations
	operations, err := client.
		DeviceManagement().
		ResourceOperations().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", readPermissions)
		return nil, err
	}

	// Create a map of valid operation IDs for quick lookup
	validOperations := make(map[string]bool)
	for _, op := range operations.GetValue() {
		if op.GetId() != nil {
			validOperations[*op.GetId()] = true
		}
	}

	// Filter out invalid permissions and log warnings
	validPermissions := []string{}
	for _, permission := range permissions {
		if validOperations[permission] {
			validPermissions = append(validPermissions, permission)
		} else {
			// Check if it's a syntax issue (e.g., missing prefix)
			if !strings.HasPrefix(permission, "Microsoft.Intune_") {
				correctedPermission := "Microsoft.Intune_" + permission
				if validOperations[correctedPermission] {
					tflog.Warn(ctx, fmt.Sprintf("Permission '%s' was missing 'Microsoft.Intune_' prefix, corrected to '%s'",
						permission, correctedPermission))
					validPermissions = append(validPermissions, correctedPermission)
					continue
				}
			}

			tflog.Warn(ctx, fmt.Sprintf("Permission '%s' is not a valid resource operation ID and will be ignored", permission))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d of %d permissions", len(validPermissions), len(permissions)))
	return validPermissions, nil
}
