package graphBetaWindowsQualityUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps the remote Windows Quality Update Profile state to the data source model
func MapRemoteStateToDataSource(ctx context.Context, data *WindowsQualityUpdateProfileDataSourceModel, remoteResource graphmodels.WindowsQualityUpdateProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote Windows Driver Quality Profile to data source model", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote Windows Driver Update Profile to data source model")
}
