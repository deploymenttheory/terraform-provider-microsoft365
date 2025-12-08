package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs an AppRoleAssignment object for creating a new app role assignment
func constructResource(ctx context.Context, data *ServicePrincipalAppRoleAssignedToResourceModel) (graphmodels.AppRoleAssignmentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAppRoleAssignment()

	if err := convert.FrameworkToGraphUUID(data.TargetServicePrincipalObjectID, requestBody.SetPrincipalId); err != nil {
		return nil, fmt.Errorf("failed to set target service principal object ID: %v", err)
	}

	if err := convert.FrameworkToGraphUUID(data.ResourceObjectID, requestBody.SetResourceId); err != nil {
		return nil, fmt.Errorf("failed to set resource object ID: %v", err)
	}

	if err := convert.FrameworkToGraphUUID(data.AppRoleID, requestBody.SetAppRoleId); err != nil {
		return nil, fmt.Errorf("failed to set app role ID: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
