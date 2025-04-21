package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs a DeviceAndAppManagementRoleAssignment from the Terraform resource model
func constructAssignment(ctx context.Context, data *RoleDefinitionResourceModel) (graphmodels.DeviceAndAppManagementRoleAssignmentable, error) {
	tflog.Debug(ctx, "Constructing role assignment")

	requestBody := graphmodels.NewDeviceAndAppManagementRoleAssignment()

	// Set basic properties
	constructors.SetStringProperty(data.Assignments.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Assignments.Description, requestBody.SetDescription)

	// Set members using the helper function
	if !data.Assignments.ScopeMembers.IsNull() && !data.Assignments.ScopeMembers.IsUnknown() {
		if err := constructors.SetStringSet(ctx, data.Assignments.ScopeMembers, requestBody.SetMembers); err != nil {
			return nil, fmt.Errorf("failed to set members: %v", err)
		}
	}

	// Set resource scopes using the helper function
	if !data.Assignments.ResourceScopes.IsNull() && !data.Assignments.ResourceScopes.IsUnknown() {
		if err := constructors.SetStringSet(ctx, data.Assignments.ResourceScopes, requestBody.SetResourceScopes); err != nil {
			return nil, fmt.Errorf("failed to set resource scopes: %v", err)
		}
	}

	// Set scope type based on schema value
	if !data.Assignments.ScopeType.IsNull() && !data.Assignments.ScopeType.IsUnknown() {
		switch data.Assignments.ScopeType.ValueString() {
		case "AllDevicesAndLicensedUsers":
			scopeType := graphmodels.ALLDEVICESANDLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		case "AllLicensedUsers":
			scopeType := graphmodels.ALLLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		case "AllDevices":
			scopeType := graphmodels.ALLDEVICES_ROLEASSIGNMENTSCOPETYPE
			requestBody.SetScopeType(&scopeType)
		default:
			return nil, fmt.Errorf("invalid scope type provided: %s", data.Assignments.ScopeType.ValueString())
		}
	}

	// Reference the created Role Definition via OData binding
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		additionalData := map[string]interface{}{
			"roleDefinition@odata.bind": fmt.Sprintf(
				"https://graph.microsoft.com/beta/deviceManagement/roleDefinitions('%s')",
				data.ID.ValueString(),
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
