package graphbetabrowsersite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *BrowserSiteResourceModel) (models.BrowserSiteable, error) {
	tflog.Debug(ctx, "Constructing BrowserSite resource")

	site := models.NewBrowserSite()

	if !data.AllowRedirect.IsNull() && !data.AllowRedirect.IsUnknown() {
		allowRedirect := data.AllowRedirect.ValueBool()
		site.SetAllowRedirect(&allowRedirect)
	}

	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		comment := data.Comment.ValueString()
		site.SetComment(&comment)
	}

	if !data.CompatibilityMode.IsNull() && !data.CompatibilityMode.IsUnknown() {
		compatibilityModeStr := data.CompatibilityMode.ValueString()
		compatibilityModeAny, err := models.ParseBrowserSiteCompatibilityMode(compatibilityModeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing CompatibilityMode: %v", err)
		}
		if compatibilityModeAny != nil {
			compatibilityMode, ok := compatibilityModeAny.(*models.BrowserSiteCompatibilityMode)
			if !ok {
				return nil, fmt.Errorf("unexpected type for CompatibilityMode: %T", compatibilityModeAny)
			}
			site.SetCompatibilityMode(compatibilityMode)
		}
	}

	if !data.MergeType.IsNull() && !data.MergeType.IsUnknown() {
		mergeTypeStr := data.MergeType.ValueString()
		mergeTypeAny, err := models.ParseBrowserSiteMergeType(mergeTypeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing MergeType: %v", err)
		}
		if mergeTypeAny != nil {
			mergeType, ok := mergeTypeAny.(*models.BrowserSiteMergeType)
			if !ok {
				return nil, fmt.Errorf("unexpected type for MergeType: %T", mergeTypeAny)
			}
			site.SetMergeType(mergeType)
		}
	}

	if !data.TargetEnvironment.IsNull() && !data.TargetEnvironment.IsUnknown() {
		targetEnvironmentStr := data.TargetEnvironment.ValueString()
		targetEnvironmentAny, err := models.ParseBrowserSiteTargetEnvironment(targetEnvironmentStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing TargetEnvironment: %v", err)
		}
		if targetEnvironmentAny != nil {
			targetEnvironment, ok := targetEnvironmentAny.(*models.BrowserSiteTargetEnvironment)
			if !ok {
				return nil, fmt.Errorf("unexpected type for TargetEnvironment: %T", targetEnvironmentAny)
			}
			site.SetTargetEnvironment(targetEnvironment)
		}
	}

	if !data.WebUrl.IsNull() && !data.WebUrl.IsUnknown() {
		webUrl := data.WebUrl.ValueString()
		site.SetWebUrl(&webUrl)
	}

	debugPrintRequestBody(ctx, site)

	return site, nil
}

func debugPrintRequestBody(ctx context.Context, site models.BrowserSiteable) {
	requestMap := map[string]interface{}{
		"allowRedirect":     site.GetAllowRedirect(),
		"comment":           site.GetComment(),
		"compatibilityMode": site.GetCompatibilityMode(),
		"mergeType":         site.GetMergeType(),
		"targetEnvironment": site.GetTargetEnvironment(),
		"webUrl":            site.GetWebUrl(),
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
