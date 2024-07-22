package assignmentFilter

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	if !data.RoleScopeTags.IsNull() && len(data.RoleScopeTags.Elements()) > 0 {
		roleScopeTags := make([]string, len(data.RoleScopeTags.Elements()))
		for i, tag := range data.RoleScopeTags.Elements() {
			roleScopeTags[i] = tag.(types.String).ValueString()
		}
		requestBody.SetRoleScopeTags(roleScopeTags)
	}

	// Set Payloads
	payloads, err := convertPayloads(data.Payloads)
	if err != nil {
		return nil, err
	}
	if payloads != nil {
		requestBody.SetPayloads(payloads)
	}

	return requestBody, nil
}

func convertPayloads(payloads types.List) ([]models.PayloadByFilterable, error) {
	if payloads.IsNull() || len(payloads.Elements()) == 0 {
		return nil, nil
	}

	result := make([]models.PayloadByFilterable, len(payloads.Elements()))
	for i, elem := range payloads.Elements() {
		payloadElem := elem.(types.Object)
		payload := models.NewPayloadByFilter()

		setStringValueFromAttributes(payloadElem.Attributes(), "payload_id", payload.SetPayloadId)
		if err := setParsedValueFromAttributes(payloadElem.Attributes(), "payload_type", func(val *models.AssociatedAssignmentPayloadType) {
			payload.SetPayloadType(val)
		}, models.ParseAssociatedAssignmentPayloadType); err != nil {
			return nil, fmt.Errorf("invalid payload type: %s", err)
		}
		setStringValueFromAttributes(payloadElem.Attributes(), "group_id", payload.SetGroupId)
		if err := setParsedValueFromAttributes(payloadElem.Attributes(), "assignment_filter_type", func(val *models.AssignmentFilterType) {
			payload.SetAssignmentFilterType(val)
		}, models.ParseAssignmentFilterType); err != nil {
			return nil, fmt.Errorf("invalid assignment filter type: %s", err)
		}

		result[i] = payload
	}
	return result, nil
}

// setStringValueFromAttributes sets a string value from the given attribute map if the key exists and is not null.
// It takes a map of attributes, a key to look for, and a setter function that sets the value if found.
func setStringValueFromAttributes(attrs map[string]attr.Value, key string, setter func(*string)) {
	if v, ok := attrs[key].(types.String); ok && !v.IsNull() {
		str := v.ValueString()
		setter(&str)
	}
}

// setParsedValueFromAttributes sets a parsed value from the given attribute map if the key exists and is not null.
// It takes a map of attributes, a key to look for, a setter function to set the parsed value, and a parser function
// to convert the string value to the desired type. It returns an error if parsing fails.
func setParsedValueFromAttributes[T any](attrs map[string]attr.Value, key string, setter func(*T), parser func(string) (interface{}, error)) error {
	if v, ok := attrs[key].(types.String); ok && !v.IsNull() {
		str := v.ValueString()
		parsedValue, err := parser(str)
		if err != nil {
			return err
		}
		if parsedValue != nil {
			setter(parsedValue.(*T))
		}
	}
	return nil
}
