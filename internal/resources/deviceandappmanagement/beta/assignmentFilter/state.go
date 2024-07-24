package graphBetaAssignmentFilter

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform
func mapRemoteStateToTerraform(data *AssignmentFilterResourceModel, remoteResource models.DeviceAndAppManagementAssignmentFilterable) {
	data.DisplayName = types.StringValue(*remoteResource.GetDisplayName())
	data.Description = types.StringValue(*remoteResource.GetDescription())
	data.Platform = types.StringValue(remoteResource.GetPlatform().String())
	data.Rule = types.StringValue(*remoteResource.GetRule())
	data.AssignmentFilterManagementType = types.StringValue(remoteResource.GetAssignmentFilterManagementType().String())
	data.CreatedDateTime = types.StringValue(remoteResource.GetCreatedDateTime().String())
	data.LastModifiedDateTime = types.StringValue(remoteResource.GetLastModifiedDateTime().String())

	// Set RoleScopeTags
	roleScopeTags := remoteResource.GetRoleScopeTags()
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
	payloads := remoteResource.GetPayloads()
	if payloads != nil {
		payloadList := make([]attr.Value, len(payloads))
		for i, payload := range payloads {
			payloadType := payload.GetPayloadType().String()
			assignmentFilterType := payload.GetAssignmentFilterType().String()
			payloadMap := map[string]attr.Value{
				"payload_id":                     types.StringValue(*payload.GetPayloadId()),
				"payload_type":                   types.StringValue(payloadType),
				"group_id":                       types.StringValue(*payload.GetGroupId()),
				"assignment_remoteResource_type": types.StringValue(assignmentFilterType),
			}
			payloadList[i] = types.ObjectValueMust(map[string]attr.Type{
				"payload_id":                     types.StringType,
				"payload_type":                   types.StringType,
				"group_id":                       types.StringType,
				"assignment_remoteResource_type": types.StringType,
			}, payloadMap)
		}
		payloadsList := types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":                     types.StringType,
				"payload_type":                   types.StringType,
				"group_id":                       types.StringType,
				"assignment_remoteResource_type": types.StringType,
			},
		}, payloadList)
		data.Payloads = payloadsList
	} else {
		payloadsList := types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":                     types.StringType,
				"payload_type":                   types.StringType,
				"group_id":                       types.StringType,
				"assignment_remoteResource_type": types.StringType,
			},
		}, []attr.Value{})
		data.Payloads = payloadsList
	}

}
