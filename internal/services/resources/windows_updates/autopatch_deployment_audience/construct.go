package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, plan *WindowsUpdatesAutopatchDeploymentAudienceResourceModel) (graphmodelswindowsupdates.DeploymentAudienceable, error) {
	tflog.Debug(ctx, "Constructing deployment audience resource")

	requestBody := graphmodelswindowsupdates.NewDeploymentAudience()

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}
