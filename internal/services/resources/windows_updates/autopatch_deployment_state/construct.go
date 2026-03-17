package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructStateResource(ctx context.Context, data *WindowsUpdatesAutopatchDeploymentStateResourceModel) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()

	deploymentState := graphmodelswindowsupdates.NewDeploymentState()

	if err := convert.FrameworkToGraphEnum(
		data.RequestedValue,
		graphmodelswindowsupdates.ParseRequestedDeploymentStateValue,
		deploymentState.SetRequestedValue,
	); err != nil {
		return nil, fmt.Errorf("failed to parse requested_value: %w", err)
	}

	requestBody.SetState(deploymentState)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

func constructResetStateResource(ctx context.Context) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing reset state request for %s", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()
	deploymentState := graphmodelswindowsupdates.NewDeploymentState()

	noneStr := "none"
	noneValue, err := graphmodelswindowsupdates.ParseRequestedDeploymentStateValue(noneStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse none state value: %w", err)
	}
	deploymentState.SetRequestedValue(noneValue.(*graphmodelswindowsupdates.RequestedDeploymentStateValue))
	requestBody.SetState(deploymentState)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing reset state request for %s", ResourceName))
	return requestBody, nil
}
