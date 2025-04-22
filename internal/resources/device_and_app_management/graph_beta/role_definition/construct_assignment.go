package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs a DeviceAndAppManagementRoleAssignment from the Terraform assignment model
func constructAssignment(ctx context.Context, roleDefinitionID string, isBuiltInRole bool, builtInRoleName string, data *sharedmodels.RoleAssignmentResourceModel) (graphmodels.DeviceAndAppManagementRoleAssignmentable, error) {
	tflog.Debug(ctx, "Constructing role assignment")

	requestBody := graphmodels.NewDeviceAndAppManagementRoleAssignment()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	if !data.ScopeMembers.IsNull() && !data.ScopeMembers.IsUnknown() {
		if err := constructors.SetStringSet(ctx, data.ScopeMembers, requestBody.SetMembers); err != nil {
			return nil, fmt.Errorf("failed to set members: %v", err)
		}
	}

	if !data.ResourceScopes.IsNull() && !data.ResourceScopes.IsUnknown() {
		if err := constructors.SetStringSet(ctx, data.ResourceScopes, requestBody.SetResourceScopes); err != nil {
			return nil, fmt.Errorf("failed to set resource scopes: %v", err)
		}
	}

	if !data.ScopeType.IsNull() && !data.ScopeType.IsUnknown() {
		switch data.ScopeType.ValueString() {
		case "allDevicesAndLicensedUsers":
			scopeType := graphmodels.ALLDEVICESANDLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		case "allLicensedUsers":
			scopeType := graphmodels.ALLLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		case "allDevices":
			scopeType := graphmodels.ALLDEVICES_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		case "resourceScope":
			scopeType := graphmodels.RESOURCESCOPE_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		default:
			return nil, fmt.Errorf("invalid scope type provided: %s", data.ScopeType.ValueString())
		}
	}

	// Determine the role definition ID to use
	definitionID := roleDefinitionID
	if isBuiltInRole && builtInRoleName != "" {
		// Look up the UUID for the built-in role name
		if uuid, exists := BuiltInIntuneRoleDefinitions[builtInRoleName]; exists {
			definitionID = uuid
			tflog.Debug(ctx, fmt.Sprintf("Using built-in role UUID %s for role name %s", uuid, builtInRoleName))
		} else {
			return nil, fmt.Errorf("unknown built-in role name: %s", builtInRoleName)
		}
	}

	if definitionID != "" {
		additionalData := map[string]interface{}{
			"roleDefinition@odata.bind": fmt.Sprintf(
				"https://graph.microsoft.com/beta/deviceManagement/roleDefinitions('%s')",
				definitionID,
			),
		}
		requestBody.SetAdditionalData(additionalData)
	} else {
		return nil, fmt.Errorf("role definition ID is required for assignment binding")
	}

	if err := constructors.DebugLogGraphObject(ctx, "Role Assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing role assignment")
	return requestBody, nil
}
