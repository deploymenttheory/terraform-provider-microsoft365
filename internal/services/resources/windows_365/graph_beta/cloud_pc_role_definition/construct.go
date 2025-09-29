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

// constructResource constructs a UnifiedRoleDefinition resource using data from the Terraform model.
func constructResource(ctx context.Context, data *RoleDefinitionResourceModel, client *msgraphbetasdk.GraphServiceClient) (graphmodels.UnifiedRoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewUnifiedRoleDefinition()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if len(data.RolePermissions) > 0 {
		unifiedRolePermissions := []graphmodels.UnifiedRolePermissionable{}

		for _, perm := range data.RolePermissions {
			if !perm.AllowedResourceActions.IsNull() {
				allowedResourceActions := []string{}
				for _, a := range perm.AllowedResourceActions.Elements() {
					if actionStr, ok := a.(types.String); ok && !actionStr.IsNull() {
						allowedResourceActions = append(allowedResourceActions, actionStr.ValueString())
					}
				}

				displayName := data.DisplayName.ValueString()
				validatedPermissions, err := validateRequest(ctx, allowedResourceActions, client, displayName)
				if err != nil {
					return nil, fmt.Errorf("failed to validate role permissions: %s", err)
				}

				unifiedRolePermission := graphmodels.NewUnifiedRolePermission()
				unifiedRolePermission.SetAllowedResourceActions(validatedPermissions)
				unifiedRolePermissions = append(unifiedRolePermissions, unifiedRolePermission)
			}
		}

		requestBody.SetRolePermissions(unifiedRolePermissions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
