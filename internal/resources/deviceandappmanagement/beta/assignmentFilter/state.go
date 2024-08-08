package graphBetaAssignmentFilter

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the AssignmentFilterResourceModel with data from the DeviceAndAppManagementAssignmentFilterable.
func mapRemoteStateToTerraform(ctx context.Context, data *AssignmentFilterResourceModel, remoteResource models.DeviceAndAppManagementAssignmentFilterable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = types.StringValue(helpers.StringPtrToString(remoteResource.GetId()))

	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetDisplayName()))

	data.Description = types.StringValue(helpers.StringPtrToString(remoteResource.GetDescription()))

	if platform := remoteResource.GetPlatform(); platform != nil {
		data.Platform = types.StringValue(platform.String())
	} else {
		data.Platform = types.StringNull()
	}

	data.Rule = types.StringValue(helpers.StringPtrToString(remoteResource.GetRule()))

	if managementType := remoteResource.GetAssignmentFilterManagementType(); managementType != nil {
		data.AssignmentFilterManagementType = types.StringValue(managementType.String())
	} else {
		data.AssignmentFilterManagementType = types.StringNull()
	}

	data.CreatedDateTime = helpers.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = helpers.TimeToString(remoteResource.GetLastModifiedDateTime())

	roleScopeTags := remoteResource.GetRoleScopeTags()

	filteredRoleScopeTags := make([]string, 0)
	for _, tag := range roleScopeTags {
		if tag != "0" { // Ignore the "0" value
			filteredRoleScopeTags = append(filteredRoleScopeTags, tag)
		}
	}

	if len(filteredRoleScopeTags) == 0 {
		data.RoleScopeTags = types.ListValueMust(types.StringType, []attr.Value{})
	} else {
		data.RoleScopeTags = types.ListValueMust(types.StringType, roleScopeTagsToValueSlice(filteredRoleScopeTags))
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform")
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
