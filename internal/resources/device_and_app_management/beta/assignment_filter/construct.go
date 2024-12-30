package graphBetaAssignmentFilter

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *AssignmentFilterResourceModel) (*graphmodels.DeviceAndAppManagementAssignmentFilter, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceAndAppManagementAssignmentFilter()

	construct.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)

	construct.SetStringProperty(data.Description, requestBody.SetDescription)

	if err := construct.SetEnumProperty(data.Platform, graphmodels.ParseDevicePlatformType, requestBody.SetPlatform); err != nil {
		return nil, fmt.Errorf("invalid device platform type: %s", err)
	}

	construct.SetStringProperty(data.Rule, requestBody.SetRule)

	if err := construct.SetEnumProperty(data.AssignmentFilterManagementType, graphmodels.ParseAssignmentFilterManagementType, requestBody.SetAssignmentFilterManagementType); err != nil {
		return nil, fmt.Errorf("invalid assignment filter management type: %s", err)
	}

	if err := construct.SetStringList(ctx, data.RoleScopeTags, requestBody.SetRoleScopeTags); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
