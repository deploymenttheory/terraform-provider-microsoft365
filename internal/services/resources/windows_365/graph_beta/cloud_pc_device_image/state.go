package graphBetaCloudPcDeviceImage

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcDeviceImageResourceModel, remoteResource models.CloudPcDeviceImageable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.ErrorCode = convert.GraphToFrameworkEnum(remoteResource.GetErrorCode())
	data.ExpirationDate = convert.GraphToFrameworkDateOnly(remoteResource.GetExpirationDate())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.OperatingSystem = convert.GraphToFrameworkString(remoteResource.GetOperatingSystem())
	data.OsBuildNumber = convert.GraphToFrameworkString(remoteResource.GetOsBuildNumber())
	data.OsStatus = convert.GraphToFrameworkEnum(remoteResource.GetOsStatus())
	data.OsVersionNumber = convert.GraphToFrameworkString(remoteResource.GetOsVersionNumber())
	data.SourceImageResourceId = convert.GraphToFrameworkString(remoteResource.GetSourceImageResourceId())
	data.Status = convert.GraphToFrameworkEnum(remoteResource.GetStatus())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
