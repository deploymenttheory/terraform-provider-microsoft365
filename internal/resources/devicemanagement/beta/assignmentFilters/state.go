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

}
