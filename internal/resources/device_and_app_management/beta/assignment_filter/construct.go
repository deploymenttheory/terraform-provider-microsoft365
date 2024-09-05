package graphBetaAssignmentFilter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *AssignmentFilterResourceModel) (*models.DeviceAndAppManagementAssignmentFilter, error) {
	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()

	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.Platform.IsNull() {
		platformStr := data.Platform.ValueString()
		platform, err := models.ParseDevicePlatformType(platformStr)
		if err != nil {
			return nil, fmt.Errorf("invalid platform: %s", err)
		}
		if platform != nil {
			requestBody.SetPlatform(platform.(*models.DevicePlatformType))
		}
	}

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	if !data.AssignmentFilterManagementType.IsNull() {
		assignmentFilterManagementTypeStr := data.AssignmentFilterManagementType.ValueString()
		assignmentFilterManagementType, err := models.ParseAssignmentFilterManagementType(assignmentFilterManagementTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid assignment filter management type: %s", err)
		}
		if assignmentFilterManagementType != nil {
			requestBody.SetAssignmentFilterManagementType(assignmentFilterManagementType.(*models.AssignmentFilterManagementType))
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

	requestBodyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed assignment filter resource:\n"+string(requestBodyJSON))

	return requestBody, nil
}
