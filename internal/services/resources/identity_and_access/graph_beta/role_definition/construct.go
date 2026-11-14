package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a UnifiedRoleDefinition that can be serialized by Kiota.
// templateId is only sent on create (the planned ID is still unknown): it becomes the role ID and cannot change afterwards.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleDefinitionResourceModel) (graphmodels.UnifiedRoleDefinitionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	isCreate := data.ID.IsNull() || data.ID.IsUnknown()
	requestBody, err := constructRequestBody(ctx, data, isCreate)
	if err != nil {
		return nil, err
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the role definition request body from the model without calling Graph
func constructRequestBody(ctx context.Context, data *RoleDefinitionResourceModel, isCreate bool) (graphmodels.UnifiedRoleDefinitionable, error) {
	requestBody := graphmodels.NewUnifiedRoleDefinition()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsEnabled, requestBody.SetIsEnabled)
	convert.FrameworkToGraphString(data.Version, requestBody.SetVersion)

	if isCreate && !data.TemplateId.IsUnknown() {
		convert.FrameworkToGraphString(data.TemplateId, requestBody.SetTemplateId)
	}

	// PATCH keeps properties that are left out, so clear description and version with an explicit null when they are removed from config.
	// Microsoft Entra rejects the empty string for both (InvalidLength), and Kiota omits nil strings, so the nulls go through additional data.
	if !isCreate {
		explicitNulls := map[string]any{}
		if data.Description.IsNull() {
			explicitNulls["description"] = nil
		}
		if data.Version.IsNull() {
			explicitNulls["version"] = nil
		}
		if len(explicitNulls) > 0 {
			requestBody.SetAdditionalData(explicitNulls)
		}
	}

	rolePermissions, err := constructRolePermissions(ctx, data.RolePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to construct role permissions: %w", err)
	}
	requestBody.SetRolePermissions(rolePermissions)

	return requestBody, nil
}

// constructRolePermissions builds the rolePermissions collection using SDK models
func constructRolePermissions(ctx context.Context, data []RoleDefinitionPermission) ([]graphmodels.UnifiedRolePermissionable, error) {
	permissions := make([]graphmodels.UnifiedRolePermissionable, 0, len(data))

	for i, permission := range data {
		rolePermission := graphmodels.NewUnifiedRolePermission()
		if err := convert.FrameworkToGraphStringSet(ctx, permission.AllowedResourceActions, rolePermission.SetAllowedResourceActions); err != nil {
			return nil, fmt.Errorf("failed to convert allowed resource actions for role permission %d: %w", i, err)
		}
		permissions = append(permissions, rolePermission)
	}

	return permissions, nil
}
