package graphBetaCloudPcDeviceImages

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateDatasource maps a Graph API response to the data source model
func StateDatasource(ctx context.Context, deviceImage graphmodels.CloudPcDeviceImageable) CloudPcDeviceImageItem {
	if deviceImage == nil {
		tflog.Debug(ctx, "Device image is nil, returning empty model")
		return CloudPcDeviceImageItem{}
	}

	tflog.Debug(ctx, "Mapping device image to state", map[string]any{
		"id": deviceImage.GetId(),
	})

	return CloudPcDeviceImageItem{
		ID:                    convert.GraphToFrameworkString(deviceImage.GetId()),
		DisplayName:           convert.GraphToFrameworkString(deviceImage.GetDisplayName()),
		ExpirationDate:        convert.GraphToFrameworkEnum(deviceImage.GetExpirationDate()),
		OSBuildNumber:         convert.GraphToFrameworkString(deviceImage.GetOsBuildNumber()),
		OSStatus:              convert.GraphToFrameworkEnum(deviceImage.GetOsStatus()),
		OperatingSystem:       convert.GraphToFrameworkString(deviceImage.GetOperatingSystem()),
		Version:               convert.GraphToFrameworkString(deviceImage.GetVersion()),
		SourceImageResourceID: convert.GraphToFrameworkString(deviceImage.GetSourceImageResourceId()),
		LastModifiedDateTime:  convert.GraphToFrameworkTime(deviceImage.GetLastModifiedDateTime()),
		Status:                convert.GraphToFrameworkEnum(deviceImage.GetStatus()),
		StatusDetails:         convert.GraphToFrameworkEnum(deviceImage.GetStatusDetails()),
		ErrorCode:             convert.GraphToFrameworkEnum(deviceImage.GetErrorCode()),
		OSVersionNumber:       convert.GraphToFrameworkString(deviceImage.GetOsVersionNumber()),
	}
}
