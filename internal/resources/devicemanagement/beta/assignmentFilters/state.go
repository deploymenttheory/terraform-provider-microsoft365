package assignmentFilter

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func setTerraformState(data *AssignmentFilterResourceModel, filter models.DeviceAndAppManagementAssignmentFilterable, resp *resource.ReadResponse) {
	data.DisplayName = types.StringValue(*filter.GetDisplayName())
	data.Description = types.StringValue(*filter.GetDescription())
	data.Platform = types.StringValue(filter.GetPlatform().String())
	data.Rule = types.StringValue(*filter.GetRule())
	data.AssignmentFilterManagementType = types.StringValue(filter.GetAssignmentFilterManagementType().String())
	data.CreatedDateTime = types.StringValue(filter.GetCreatedDateTime().String())
	data.LastModifiedDateTime = types.StringValue(filter.GetLastModifiedDateTime().String())

	// Set RoleScopeTags
	roleScopeTags := filter.GetRoleScopeTags()
	if roleScopeTags != nil {
		tagList := make([]attr.Value, len(roleScopeTags))
		for i, tag := range roleScopeTags {
			tagList[i] = types.StringValue(tag)
		}
		roleScopeTagsList := types.ListValueMust(types.StringType, tagList)
		data.RoleScopeTags = roleScopeTagsList
	} else {
		roleScopeTagsList := types.ListValueMust(types.StringType, []attr.Value{})
		data.RoleScopeTags = roleScopeTagsList
	}

	// Set Payloads
	payloads := filter.GetPayloads()
	if payloads != nil {
		payloadList := make([]attr.Value, len(payloads))
		for i, payload := range payloads {
			payloadType := string(*payload.GetPayloadType())
			assignmentFilterType := string(*payload.GetAssignmentFilterType())
			payloadMap := map[string]attr.Value{
				"payload_id":             types.StringValue(*payload.GetPayloadId()),
				"payload_type":           types.StringValue(payloadType),
				"group_id":               types.StringValue(*payload.GetGroupId()),
				"assignment_filter_type": types.StringValue(assignmentFilterType),
			}
			payloadList[i] = types.ObjectValueMust(map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			}, payloadMap)
		}
		payloadsList := types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			},
		}, payloadList)
		data.Payloads = payloadsList
	} else {
		payloadsList := types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			},
		}, []attr.Value{})
		data.Payloads = payloadsList
	}

}
