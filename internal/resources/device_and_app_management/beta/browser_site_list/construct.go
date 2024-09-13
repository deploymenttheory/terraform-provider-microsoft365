package graphbetabrowsersite

import (
	"context"
	"encoding/json"

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

	debugPrintRequestBody(ctx, siteList)

	return siteList, nil
}

func debugPrintRequestBody(ctx context.Context, siteList models.BrowserSiteListable) {
	requestMap := map[string]interface{}{
		"description": siteList.GetDescription(),
		"displayName": siteList.GetDisplayName(),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}
