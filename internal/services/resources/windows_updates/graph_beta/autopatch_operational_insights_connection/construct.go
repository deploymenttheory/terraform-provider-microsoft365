package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchOperationalInsightsConnectionResourceModel) (graphmodelswindowsupdates.OperationalInsightsConnectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewOperationalInsightsConnection()

	convert.FrameworkToGraphString(data.AzureResourceGroupName, requestBody.SetAzureResourceGroupName)
	convert.FrameworkToGraphString(data.AzureSubscriptionId, requestBody.SetAzureSubscriptionId)
	convert.FrameworkToGraphString(data.WorkspaceName, requestBody.SetWorkspaceName)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}
