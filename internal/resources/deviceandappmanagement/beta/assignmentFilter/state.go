package graphBetaAssignmentFilter

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform
func mapRemoteStateToTerraform(data *AssignmentFilterResourceModel, remoteResource models.DeviceAndAppManagementAssignmentFilterable) {
	if displayName := remoteResource.GetDisplayName(); displayName != nil {
		data.DisplayName = types.StringValue(*displayName)
	} else {
		data.DisplayName = types.StringNull()
	}

	if description := remoteResource.GetDescription(); description != nil {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	if platform := remoteResource.GetPlatform(); platform != nil {
		data.Platform = types.StringValue(platform.String())
	} else {
		data.Platform = types.StringNull()
	}

	if rule := remoteResource.GetRule(); rule != nil {
		data.Rule = types.StringValue(*rule)
	} else {
		data.Rule = types.StringNull()
	}

	if managementType := remoteResource.GetAssignmentFilterManagementType(); managementType != nil {
		data.AssignmentFilterManagementType = types.StringValue(managementType.String())
	} else {
		data.AssignmentFilterManagementType = types.StringNull()
	}

	if createdDateTime := remoteResource.GetCreatedDateTime(); createdDateTime != nil {
		data.CreatedDateTime = types.StringValue(createdDateTime.Format(helpers.TimeFormatRFC3339))
	} else {
		data.CreatedDateTime = types.StringNull()
	}

	if lastModifiedDateTime := remoteResource.GetLastModifiedDateTime(); lastModifiedDateTime != nil {
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime.Format(helpers.TimeFormatRFC3339))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	// Set RoleScopeTags
	if roleScopeTags := remoteResource.GetRoleScopeTags(); roleScopeTags != nil {
		tagList := make([]attr.Value, len(roleScopeTags))
		for i, tag := range roleScopeTags {
			tagList[i] = types.StringValue(tag)
		}
		data.RoleScopeTags = types.ListValueMust(types.StringType, tagList)
	} else {
		data.RoleScopeTags = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Set Payloads
	if payloads := remoteResource.GetPayloads(); payloads != nil {
		payloadList := make([]attr.Value, 0, len(payloads))
		for _, payload := range payloads {
			payloadMap := map[string]attr.Value{
				"payload_id":             types.StringValue(common.SafeDeref(payload.GetPayloadId())),
				"payload_type":           types.StringValue(common.SafeEnumString(payload.GetPayloadType())),
				"group_id":               types.StringValue(common.SafeDeref(payload.GetGroupId())),
				"assignment_filter_type": types.StringValue(common.SafeEnumString(payload.GetAssignmentFilterType())),
			}
			payloadObj, diags := types.ObjectValue(map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			}, payloadMap)
			if diags.HasError() {
				// Handle or log the error
				continue
			}
			payloadList = append(payloadList, payloadObj)
		}
		data.Payloads, _ = types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			},
		}, payloadList)
	} else {
		data.Payloads = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id":             types.StringType,
				"payload_type":           types.StringType,
				"group_id":               types.StringType,
				"assignment_filter_type": types.StringType,
			},
		})
	}
}
