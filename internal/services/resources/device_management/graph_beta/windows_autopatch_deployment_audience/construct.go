package graphBetaWindowsAutopatchDeploymentAudience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, plan *WindowsAutopatchDeploymentAudienceResourceModel) (graphmodelswindowsupdates.DeploymentAudienceable, error) {
	tflog.Debug(ctx, "Constructing deployment audience resource")

	requestBody := graphmodelswindowsupdates.NewDeploymentAudience()

	tflog.Debug(ctx, "Finished constructing deployment audience resource")
	return requestBody, nil
}
