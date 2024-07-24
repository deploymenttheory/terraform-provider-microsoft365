package graphBetaAssignmentFilter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
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

	if !data.RoleScopeTags.IsNull() && len(data.RoleScopeTags.Elements()) > 0 {
		roleScopeTags := make([]string, len(data.RoleScopeTags.Elements()))
		for i, tag := range data.RoleScopeTags.Elements() {
			roleScopeTags[i] = tag.(types.String).ValueString()
		}
		requestBody.SetRoleScopeTags(roleScopeTags)
	}

	payloads, err := convertPayloads(data.Payloads)
	if err != nil {
		return nil, err
	}
	if payloads != nil {
		requestBody.SetPayloads(payloads)
	}

	tflog.Debug(ctx, "Constructed assignment filter resource", map[string]interface{}{
		"displayName":    requestBody.GetDisplayName(),
		"description":    requestBody.GetDescription(),
		"platform":       requestBody.GetPlatform(),
		"rule":           requestBody.GetRule(),
		"managementType": requestBody.GetAssignmentFilterManagementType(),
		"roleScopeTags":  requestBody.GetRoleScopeTags(),
		"payloads":       payloads,
	})

	return requestBody, nil
}

// convertPayloads
func convertPayloads(payloads types.List) ([]models.PayloadByFilterable, error) {
	if payloads.IsNull() || len(payloads.Elements()) == 0 {
		return nil, nil
	}

	result := make([]models.PayloadByFilterable, len(payloads.Elements()))
	for i, elem := range payloads.Elements() {
		payloadElem := elem.(types.Object)
		payload := models.NewPayloadByFilter()

		common.SetStringValueFromAttributes(payloadElem.Attributes(), "payload_id", payload.SetPayloadId)
		if err := common.SetParsedValueFromAttributes(payloadElem.Attributes(), "payload_type", func(val *models.AssociatedAssignmentPayloadType) {
			payload.SetPayloadType(val)
		}, models.ParseAssociatedAssignmentPayloadType); err != nil {
			return nil, fmt.Errorf("invalid payload type: %s", err)
		}
		common.SetStringValueFromAttributes(payloadElem.Attributes(), "group_id", payload.SetGroupId)
		if err := common.SetParsedValueFromAttributes(payloadElem.Attributes(), "assignment_filter_type", func(val *models.DeviceAndAppManagementAssignmentFilterType) {
			payload.SetAssignmentFilterType(val)
		}, models.ParseAssociatedAssignmentPayloadType); err != nil {
			return nil, fmt.Errorf("invalid assignment filter type: %s", err)
		}

		result[i] = payload
	}
	return result, nil
}
