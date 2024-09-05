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

	requestBodyJSON, err := json.MarshalIndent(site, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed BrowserSite resource:\n"+string(requestBodyJSON))

	return site, nil
}

func constructHistory(ctx context.Context, history []BrowserSiteHistoryModel) ([]models.BrowserSiteHistoryable, error) {
	var constructedHistory []models.BrowserSiteHistoryable

	for _, historyItem := range history {
		newHistoryItem := models.NewBrowserSiteHistory()

		if !historyItem.AllowRedirect.IsNull() && !historyItem.AllowRedirect.IsUnknown() {
			allowRedirect := historyItem.AllowRedirect.ValueBool()
			newHistoryItem.SetAllowRedirect(&allowRedirect)
		}

		if !historyItem.Comment.IsNull() && !historyItem.Comment.IsUnknown() {
			comment := historyItem.Comment.ValueString()
			newHistoryItem.SetComment(&comment)
		}

		if !historyItem.CompatibilityMode.IsNull() && !historyItem.CompatibilityMode.IsUnknown() {
			compatibilityModeStr := historyItem.CompatibilityMode.ValueString()
			compatibilityModeAny, err := models.ParseBrowserSiteCompatibilityMode(compatibilityModeStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing CompatibilityMode in history: %v", err)
			}
			if compatibilityModeAny != nil {
				compatibilityMode, ok := compatibilityModeAny.(*models.BrowserSiteCompatibilityMode)
				if !ok {
					return nil, fmt.Errorf("unexpected type for CompatibilityMode in history: %T", compatibilityModeAny)
				}
				newHistoryItem.SetCompatibilityMode(compatibilityMode)
			}
		}

		if !historyItem.MergeType.IsNull() && !historyItem.MergeType.IsUnknown() {
			mergeTypeStr := historyItem.MergeType.ValueString()
			mergeTypeAny, err := models.ParseBrowserSiteMergeType(mergeTypeStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing MergeType in history: %v", err)
			}
			if mergeTypeAny != nil {
				mergeType, ok := mergeTypeAny.(*models.BrowserSiteMergeType)
				if !ok {
					return nil, fmt.Errorf("unexpected type for MergeType in history: %T", mergeTypeAny)
				}
				newHistoryItem.SetMergeType(mergeType)
			}
		}

		if !historyItem.TargetEnvironment.IsNull() && !historyItem.TargetEnvironment.IsUnknown() {
			targetEnvironmentStr := historyItem.TargetEnvironment.ValueString()
			targetEnvironmentAny, err := models.ParseBrowserSiteTargetEnvironment(targetEnvironmentStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing TargetEnvironment in history: %v", err)
			}
			if targetEnvironmentAny != nil {
				targetEnvironment, ok := targetEnvironmentAny.(*models.BrowserSiteTargetEnvironment)
				if !ok {
					return nil, fmt.Errorf("unexpected type for TargetEnvironment in history: %T", targetEnvironmentAny)
				}
				newHistoryItem.SetTargetEnvironment(targetEnvironment)
			}
		}

		constructedHistory = append(constructedHistory, newHistoryItem)
	}

	historyJSON, err := json.MarshalIndent(constructedHistory, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling history to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed BrowserSite history:\n"+string(historyJSON))

	return constructedHistory, nil
}
