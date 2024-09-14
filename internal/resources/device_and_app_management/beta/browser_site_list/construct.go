package graphbetabrowsersite

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *BrowserSiteListResourceModel) (models.BrowserSiteListable, error) {
	tflog.Debug(ctx, "Constructing BrowserSiteList resource")
	construct.DebugPrintStruct(ctx, "Constructed Browser Site List resource from model", data)

	siteList := models.NewBrowserSiteList()

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		siteList.SetDescription(&description)
	}

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		siteList.SetDisplayName(&displayName)
	}

	return siteList, nil
}
