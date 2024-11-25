package graphBetaAssignmentFilter

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *AssignmentFilterResourceModel) (*graphmodels.DeviceAndAppManagementAssignmentFilter, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceAndAppManagementAssignmentFilter()

	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.Platform.IsNull() {
		platformStr := data.Platform.ValueString()
		platform, err := graphmodels.ParseDevicePlatformType(platformStr)
		if err != nil {
			return nil, fmt.Errorf("invalid platform: %s", err)
		}
		if platform != nil {
			requestBody.SetPlatform(platform.(*graphmodels.DevicePlatformType))
		}
	}

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	if !data.AssignmentFilterManagementType.IsNull() {
		assignmentFilterManagementTypeStr := data.AssignmentFilterManagementType.ValueString()
		assignmentFilterManagementType, err := graphmodels.ParseAssignmentFilterManagementType(assignmentFilterManagementTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid assignment filter management type: %s", err)
		}
		if assignmentFilterManagementType != nil {
			requestBody.SetAssignmentFilterManagementType(assignmentFilterManagementType.(*graphmodels.AssignmentFilterManagementType))
		}
	}

	roleScopeTags := make([]string, 0)
	if !data.RoleScopeTags.IsNull() {
		for _, tag := range data.RoleScopeTags.Elements() {
			tagValue := tag.(types.String).ValueString()
			if tagValue != "0" {
				roleScopeTags = append(roleScopeTags, tagValue)
			}
		}
	}
	requestBody.SetRoleScopeTags(roleScopeTags)

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
