package graphBetaGroupAppRoleAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs an AppRoleAssignment object for creating a new app role assignment
func constructResource(ctx context.Context, data *GroupAppRoleAssignmentResourceModel) (graphmodels.AppRoleAssignmentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAppRoleAssignment()

	convert.FrameworkToGraphUUID(data.TargetGroupID, requestBody.SetPrincipalId)
	convert.FrameworkToGraphUUID(data.ResourceObjectID, requestBody.SetResourceId)
	convert.FrameworkToGraphUUID(data.AppRoleID, requestBody.SetAppRoleId)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
