package graphBetaAssignmentFilter

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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
		"resourceId": types.StringPointerValue(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Platform = state.EnumPtrToTypeString(remoteResource.GetPlatform())
	data.Rule = types.StringPointerValue(remoteResource.GetRule())
	data.AssignmentFilterManagementType = state.EnumPtrToTypeString(remoteResource.GetAssignmentFilterManagementType())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTags = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTags())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
