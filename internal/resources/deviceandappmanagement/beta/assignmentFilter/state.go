package graphBetaAssignmentFilter

import (
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the AssignmentFilterResourceModel with data from the DeviceAndAppManagementAssignmentFilterable.
func mapRemoteStateToTerraform(data *AssignmentFilterResourceModel, remoteResource models.DeviceAndAppManagementAssignmentFilterable) {
	if remoteResource == nil {
		return
	}

	if id := remoteResource.GetId(); id != nil {
		data.ID = types.StringValue(*id)
	}

	if displayName := remoteResource.GetDisplayName(); displayName != nil {
		data.DisplayName = types.StringValue(*displayName)
	}

	if description := remoteResource.GetDescription(); description != nil {
		data.Description = types.StringValue(*description)
	}

	if platform := remoteResource.GetPlatform(); platform != nil {
		data.Platform = types.StringValue(platform.String())
	}

	if rule := remoteResource.GetRule(); rule != nil {
		data.Rule = types.StringValue(*rule)
	}

	if managementType := remoteResource.GetAssignmentFilterManagementType(); managementType != nil {
		data.AssignmentFilterManagementType = types.StringValue((*managementType).String())
	}

	if createdDateTime := remoteResource.GetCreatedDateTime(); createdDateTime != nil {
		data.CreatedDateTime = types.StringValue(createdDateTime.Format(time.RFC3339))
	}

	if lastModifiedDateTime := remoteResource.GetLastModifiedDateTime(); lastModifiedDateTime != nil {
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime.Format(time.RFC3339))
	}

	if roleScopeTags := remoteResource.GetRoleScopeTags(); roleScopeTags != nil {
		data.RoleScopeTags = types.ListValueMust(types.StringType, roleScopeTagsToValueSlice(roleScopeTags))
	}

	if payloads := remoteResource.GetPayloads(); payloads != nil {
		data.Payloads = types.ListValueMust(types.ObjectType{AttrTypes: payloadAttributeTypes()}, payloadsToValueSlice(payloads))
	}
}

// roleScopeTagsToValueSlice converts a slice of role scope tag strings to a slice of Terraform attr.Value.
// This is used to populate the RoleScopeTags field in the Terraform resource model.
func roleScopeTagsToValueSlice(roleScopeTags []string) []attr.Value {
	values := make([]attr.Value, len(roleScopeTags))
	for i, tag := range roleScopeTags {
		values[i] = types.StringValue(tag)
	}
	return values
}

// payloadAttributeTypes returns a map of attribute names to their Terraform types for the Payload object.
// This defines the structure of the Payload object in the Terraform resource model.
func payloadAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"payload_id":             types.StringType,
		"payload_type":           types.StringType,
		"group_id":               types.StringType,
		"assignment_filter_type": types.StringType,
	}
}

// payloadsToValueSlice converts a slice of PayloadByFilterable to a slice of Terraform attr.Value.
// This is used to populate the Payloads field in the Terraform resource model.
func payloadsToValueSlice(payloads []models.PayloadByFilterable) []attr.Value {
	values := make([]attr.Value, len(payloads))
	for i, payload := range payloads {
		payloadMap := map[string]attr.Value{
			"payload_id":             types.StringValue(common.StringPtrToString(payload.GetPayloadId())),
			"payload_type":           types.StringValue(string(*payload.GetPayloadType())),
			"group_id":               types.StringValue(common.StringPtrToString(payload.GetGroupId())),
			"assignment_filter_type": types.StringValue(string(*payload.GetAssignmentFilterType())),
		}
		values[i] = types.ObjectValueMust(payloadAttributeTypes(), payloadMap)
	}
	return values
}
