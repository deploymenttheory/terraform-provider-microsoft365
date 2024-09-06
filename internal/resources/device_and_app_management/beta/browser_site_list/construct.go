package graphbetabrowsersite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *BrowserSiteListResourceModel) (models.BrowserSiteListable, error) {
	tflog.Debug(ctx, "Constructing BrowserSiteList resource")

	siteList := models.NewBrowserSiteList()

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		siteList.SetDescription(&description)
	}

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		siteList.SetDisplayName(&displayName)
	}

	requestBodyJSON, err := json.MarshalIndent(siteList, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed BrowserSiteList resource:\n"+string(requestBodyJSON))

	return siteList, nil
}
