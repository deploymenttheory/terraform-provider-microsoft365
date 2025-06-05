package graphBetaDeviceEnrollmentLimitConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *DeviceEnrollmentLimitConfigurationResourceModel, remoteResource models.DeviceEnrollmentLimitConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": types.StringPointerValue(remoteResource.GetId()),
	})

	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.Priority = state.Int32PointerValue(remoteResource.GetPriority())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.Version = state.Int32PointerValue(remoteResource.GetVersion())
	data.Limit = state.Int32PointerValue(remoteResource.GetLimit())

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		data.DeviceEnrollmentConfigurationType = types.StringValue(configType.String())
	}

	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		data.RoleScopeTagIds = state.StringSliceToSet(ctx, roleScopeTagIds)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
