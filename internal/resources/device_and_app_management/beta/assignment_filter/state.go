package graphBetaAssignmentFilter

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an AssignmentFilterResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *AssignmentFilterResourceModel, remoteResource graphmodels.DeviceAndAppManagementAssignmentFilterable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.Platform = state.EnumPtrToTypeString(remoteResource.GetPlatform())
	data.Rule = types.StringValue(state.StringPtrToString(remoteResource.GetRule()))
	data.AssignmentFilterManagementType = state.EnumPtrToTypeString(remoteResource.GetAssignmentFilterManagementType())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())

	// Special handling for RoleScopeTags
	var roleScopeTagIds []attr.Value
	for _, v := range state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTags()) {
		roleScopeTagIds = append(roleScopeTagIds, v)
	}

	data.RoleScopeTags = types.ListValueMust(
		types.StringType,
		roleScopeTagIds,
	)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
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
