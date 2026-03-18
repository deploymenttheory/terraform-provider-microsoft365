package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"context"
	"fmt"

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

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
