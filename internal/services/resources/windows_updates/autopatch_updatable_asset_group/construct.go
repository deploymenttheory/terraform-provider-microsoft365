package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, _ *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel) (graphmodelswindowsupdates.UpdatableAssetGroupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewUpdatableAssetGroup()

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
