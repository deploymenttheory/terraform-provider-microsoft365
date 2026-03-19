package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, _ *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel) (graphmodelswindowsupdates.UpdatableAssetGroupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewUpdatableAssetGroup()

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}
