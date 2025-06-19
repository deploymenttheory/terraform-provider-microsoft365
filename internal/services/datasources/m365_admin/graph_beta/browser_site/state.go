package graphBetaBrowserSite

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Graph API BrowserSite to the data source model
func MapRemoteStateToDataSource(ctx context.Context, data *BrowserSiteResourceModel, remoteResource graphmodels.BrowserSiteable, browserSiteListId types.String) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to data source item", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.BrowserSiteListAssignmentID = browserSiteListId
	data.WebUrl = convert.GraphToFrameworkString(remoteResource.GetWebUrl())

	tflog.Debug(ctx, "Finished mapping remote state to data source item", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
