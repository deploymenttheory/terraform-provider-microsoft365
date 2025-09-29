package graphBetaBrowserSite

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource is the main entry point to construct the browser site resource for the Terraform provider.
func constructResource(ctx context.Context, data *BrowserSiteResourceModel) (graphmodels.BrowserSiteable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewBrowserSite()

	convert.FrameworkToGraphBool(data.AllowRedirect, requestBody.SetAllowRedirect)
	convert.FrameworkToGraphString(data.Comment, requestBody.SetComment)
	convert.FrameworkToGraphString(data.WebUrl, requestBody.SetWebUrl)

	if err := convert.FrameworkToGraphEnum(data.CompatibilityMode,
		graphmodels.ParseBrowserSiteCompatibilityMode,
		requestBody.SetCompatibilityMode); err != nil {
		return nil, fmt.Errorf("failed to set compatibility mode: %v", err)
	}

	if err := convert.FrameworkToGraphEnum(data.MergeType,
		graphmodels.ParseBrowserSiteMergeType,
		requestBody.SetMergeType); err != nil {
		return nil, fmt.Errorf("failed to set merge type: %v", err)
	}

	if err := convert.FrameworkToGraphEnum(data.TargetEnvironment,
		graphmodels.ParseBrowserSiteTargetEnvironment,
		requestBody.SetTargetEnvironment); err != nil {
		return nil, fmt.Errorf("failed to set target environment: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
