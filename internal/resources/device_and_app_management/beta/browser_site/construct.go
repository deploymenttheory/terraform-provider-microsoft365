package graphBetaBrowserSite

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs an assignment filter resource using data from the Terraform model.
func constructResource(ctx context.Context, data *BrowserSiteResourceModel) (models.BrowserSiteable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewBrowserSite()

	if !data.AllowRedirect.IsNull() && !data.AllowRedirect.IsUnknown() {
		allowRedirect := data.AllowRedirect.ValueBool()
		requestBody.SetAllowRedirect(&allowRedirect)
	}

	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		comment := data.Comment.ValueString()
		requestBody.SetComment(&comment)
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
			requestBody.SetCompatibilityMode(compatibilityMode)
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
			requestBody.SetMergeType(mergeType)
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
			requestBody.SetTargetEnvironment(targetEnvironment)
		}
	}

	if !data.WebUrl.IsNull() && !data.WebUrl.IsUnknown() {
		webUrl := data.WebUrl.ValueString()
		requestBody.SetWebUrl(&webUrl)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
