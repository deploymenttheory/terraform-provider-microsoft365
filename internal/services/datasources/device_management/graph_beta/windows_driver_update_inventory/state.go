package graphBetaWindowsDriverUpdateInventory

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps the remote Windows Driver Update Inventory state to the data source model
func MapRemoteStateToDataSource(ctx context.Context, data *WindowsDriverUpdateInventoryDataSourceModel, remoteResource graphmodels.WindowsDriverUpdateInventoryable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote Windows Driver Update Inventory to data source model", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())

	tflog.Debug(ctx, "Finished mapping remote Windows Driver Update Inventory to data source model")
}
