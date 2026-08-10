package graphBetaRoleDefinitionAssignment

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

var (
	errConvertMembersSet      = errors.New("failed to convert members set to slice")
	errConvertRoleScopeTagSet = errors.New("failed to convert role_scope_tag_ids set to slice")
	errConvertResourceScopes  = errors.New("failed to convert resource_scopes set to slice")
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(
	ctx context.Context,
	data *RoleAssignmentResourceModel,
) (graphmodels.DeviceAndAppManagementRoleAssignmentable, error) {
	requestBody := graphmodels.NewDeviceAndAppManagementRoleAssignment()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	// Set members
	if !data.Members.IsNull() && !data.Members.IsUnknown() {
		membersElements := make([]types.String, 0, len(data.Members.Elements()))
		diags := data.Members.ElementsAs(ctx, &membersElements, false)
		if diags.HasError() {
			return nil, fmt.Errorf("%w", errConvertMembersSet)
		}

		members := make([]string, len(membersElements))
		for i, member := range membersElements {
			members[i] = member.ValueString()
		}
		requestBody.SetMembers(members)
	}

	// Set role scope tag IDs
	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		roleScopeTagElements := make([]types.String, 0, len(data.RoleScopeTagIds.Elements()))
		diags := data.RoleScopeTagIds.ElementsAs(ctx, &roleScopeTagElements, false)
		if diags.HasError() {
			return nil, fmt.Errorf("%w", errConvertRoleScopeTagSet)
		}

		roleScopeTagIds := make([]string, len(roleScopeTagElements))
		for i, scopeTag := range roleScopeTagElements {
			roleScopeTagIds[i] = scopeTag.ValueString()
		}
		requestBody.SetRoleScopeTagIds(roleScopeTagIds)
	}

	// Handle scope configuration
	if len(data.ScopeConfig) > 0 {
		scopeConfig := data.ScopeConfig[0]

		switch scopeConfig.Type.ValueString() {
		case "ResourceScopes":
			// Set resource sc		opes
			if !scopeConfig.ResourceScopes.IsNull() && !scopeConfig.ResourceScopes.IsUnknown() {
				resourceScopesElements := make(
					[]types.String,
					0,
					len(scopeConfig.ResourceScopes.Elements()),
				)
				diags := scopeConfig.ResourceScopes.ElementsAs(ctx, &resourceScopesElements, false)
				if diags.HasError() {
					return nil, fmt.Errorf("%w", errConvertResourceScopes)
				}

				resourceScopes := make([]string, len(resourceScopesElements))
				for i, scope := range resourceScopesElements {
					resourceScopes[i] = scope.ValueString()
				}
				requestBody.SetResourceScopes(resourceScopes)
			}
		case "AllLicensedUsers":
			// Set scope type to all licensed users
			scopeType := graphmodels.ALLLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
			// Set empty resource scopes
			requestBody.SetResourceScopes([]string{})
		case "AllDevices":
			// Set scope type to all devices
			scopeType := graphmodels.ALLDEVICES_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
			// Set empty resource scopes
			requestBody.SetResourceScopes([]string{})
		}
	}

	// Set role definition using odata.bind
	if !data.RoleDefinitionId.IsNull() && !data.RoleDefinitionId.IsUnknown() {
		roleDefId := data.RoleDefinitionId.ValueString()
		additionalData := map[string]any{
			"roleDefinition@odata.bind": fmt.Sprintf(
				"https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions/%s",
				roleDefId,
			),
		}
		requestBody.SetAdditionalData(additionalData)
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
