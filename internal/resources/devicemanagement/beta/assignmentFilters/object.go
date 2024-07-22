package assignmentFilter

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(data *AssignmentFilterResourceModel) (*models.DeviceAndAppManagementAssignmentFilter, error) {
	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()

	// Set DisplayName
	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	// Set Description
	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	// Set Platform
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

	// Set Rule
	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	// Set AssignmentFilterManagementType
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

	// Set RoleScopeTags
	if !data.RoleScopeTags.IsNull() {
		var roleScopeTags []string
		for _, tag := range data.RoleScopeTags.Elements() {
			roleScopeTags = append(roleScopeTags, tag.(types.String).Value)
		}
		requestBody.SetRoleScopeTags(roleScopeTags)
	}

	// Set Payloads
	if !data.Payloads.IsNull() {
		var payloads []models.Payload
		for _, payloadElement := range data.Payloads.Elements() {
			payload := payloadElement.(types.Object)
			payloadID := payload.Attributes["payload_id"].(types.String).ValueString()
			payloadType := payload.Attributes["payload_type"].(types.String).ValueString()
			groupID := payload.Attributes["group_id"].(types.String).ValueString()
			assignmentFilterType := payload.Attributes["assignment_filter_type"].(types.String).ValueString()

			payloadModel := models.NewPayload()
			payloadModel.SetPayloadId(&payloadID)
			payloadModel.SetPayloadType(&payloadType)
			payloadModel.SetGroupId(&groupID)
			payloadModel.SetAssignmentFilterType(&assignmentFilterType)

			payloads = append(payloads, payloadModel)
		}
		requestBody.SetPayloads(payloads)
	}

	return requestBody, nil
}
