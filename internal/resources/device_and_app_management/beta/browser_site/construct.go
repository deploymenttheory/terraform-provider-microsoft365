package graphBetaBrowserSite

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource is the main entry point to construct the browser site resource for the Terraform provider.
func constructResource(ctx context.Context, data *BrowserSiteResourceModel) (graphmodels.BrowserSiteable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewBrowserSite()

	// Set basic properties
	constructors.SetBoolProperty(data.AllowRedirect, requestBody.SetAllowRedirect)
	constructors.SetStringProperty(data.Comment, requestBody.SetComment)
	constructors.SetStringProperty(data.WebUrl, requestBody.SetWebUrl)

	// Handle compatibility mode enum
	if err := constructors.SetEnumProperty(data.CompatibilityMode,
		graphmodels.ParseBrowserSiteCompatibilityMode,
		requestBody.SetCompatibilityMode); err != nil {
		return nil, fmt.Errorf("failed to set compatibility mode: %v", err)
	}

	// Handle merge type enum
	if err := constructors.SetEnumProperty(data.MergeType,
		graphmodels.ParseBrowserSiteMergeType,
		requestBody.SetMergeType); err != nil {
		return nil, fmt.Errorf("failed to set merge type: %v", err)
	}

	// Handle target environment enum
	if err := constructors.SetEnumProperty(data.TargetEnvironment,
		graphmodels.ParseBrowserSiteTargetEnvironment,
		requestBody.SetTargetEnvironment); err != nil {
		return nil, fmt.Errorf("failed to set target environment: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
