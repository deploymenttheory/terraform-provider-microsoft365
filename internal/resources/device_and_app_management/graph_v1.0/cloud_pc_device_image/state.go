package graphCloudPcDeviceImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcDeviceImageResourceModel, remoteResource models.CloudPcDeviceImageable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.ErrorCode = state.EnumPtrToTypeString(remoteResource.GetErrorCode())
	data.ExpirationDate = state.DateOnlyPtrToString(remoteResource.GetExpirationDate())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.OperatingSystem = types.StringPointerValue(remoteResource.GetOperatingSystem())
	data.OsBuildNumber = types.StringPointerValue(remoteResource.GetOsBuildNumber())
	data.OsStatus = state.EnumPtrToTypeString(remoteResource.GetOsStatus())
	data.SourceImageResourceId = types.StringPointerValue(remoteResource.GetSourceImageResourceId())
	data.Status = state.EnumPtrToTypeString(remoteResource.GetStatus())
	data.Version = types.StringPointerValue(remoteResource.GetVersion())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
