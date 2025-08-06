package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a RoleDefinition resource using data from the Terraform model.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleDefinitionResourceModel, resp interface{}, readPermissions []string) (graphmodels.RoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewRoleDefinition()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if len(data.RolePermissions) > 0 {
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

		validatedPermissions, err := validateRequest(ctx, client, allowedResourceActions, resp, readPermissions)
		if err != nil {
			return nil, fmt.Errorf("failed to validate role permissions: %s", err)
		}

		resourceAction.SetAllowedResourceActions(validatedPermissions)
		resourceActions := []graphmodels.ResourceActionable{resourceAction}
		rolePermission.SetResourceActions(resourceActions)
		rolePermissions := []graphmodels.RolePermissionable{rolePermission}
		requestBody.SetRolePermissions(rolePermissions)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
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
