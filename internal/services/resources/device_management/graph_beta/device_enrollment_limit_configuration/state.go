package graphBetaDeviceEnrollmentLimitConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *DeviceEnrollmentLimitConfigurationResourceModel, remoteResource models.DeviceEnrollmentLimitConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Priority = convert.GraphToFrameworkInt32(remoteResource.GetPriority())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())
	data.Limit = convert.GraphToFrameworkInt32(remoteResource.GetLimit())

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		data.DeviceEnrollmentConfigurationType = types.StringValue(configType.String())
	}

	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, roleScopeTagIds)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
