package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote role definition from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data RoleDefinitionResourceModel, remoteResource graphmodels.UnifiedRoleDefinitionable) RoleDefinitionResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = mapOptionalString(remoteResource.GetDescription())
	data.IsEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsEnabled())
	data.TemplateId = convert.GraphToFrameworkString(remoteResource.GetTemplateId())
	data.Version = mapOptionalString(remoteResource.GetVersion())
	data.IsBuiltIn = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())
	data.IsPrivileged = convert.GraphToFrameworkBool(remoteResource.GetIsPrivileged())
	data.ResourceScopes = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetResourceScopes())
	data.RolePermissions = mapRolePermissions(ctx, remoteResource.GetRolePermissions())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}

// mapRolePermissions maps the remote rolePermissions collection to the Terraform model
func mapRolePermissions(ctx context.Context, permissions []graphmodels.UnifiedRolePermissionable) []RoleDefinitionPermission {
	result := make([]RoleDefinitionPermission, 0, len(permissions))

	for _, permission := range permissions {
		if permission == nil {
			continue
		}
		result = append(result, RoleDefinitionPermission{
			AllowedResourceActions: convert.GraphToFrameworkStringSetPreserveEmpty(ctx, permission.GetAllowedResourceActions()),
		})
	}

	return result
}
