package graphBetaAssignmentFilter

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *AssignmentFilterResourceModel, isCreate bool) (*graphmodels.DeviceAndAppManagementAssignmentFilter, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceAndAppManagementAssignmentFilter()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Platform field can only be set during creation, not during updates
	if isCreate {
		if err := convert.FrameworkToGraphEnum(data.Platform, graphmodels.ParseDevicePlatformType, requestBody.SetPlatform); err != nil {
			return nil, fmt.Errorf("invalid device platform type: %s", err)
		}
	}

	convert.FrameworkToGraphString(data.Rule, requestBody.SetRule)

	if err := convert.FrameworkToGraphEnum(data.AssignmentFilterManagementType, graphmodels.ParseAssignmentFilterManagementType, requestBody.SetAssignmentFilterManagementType); err != nil {
		return nil, fmt.Errorf("invalid assignment filter management type: %s", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTags, requestBody.SetRoleScopeTags); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
