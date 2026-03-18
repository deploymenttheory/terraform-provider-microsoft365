package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentTestResource struct{}

func (r WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		groupId := state.Attributes["updatable_asset_group_id"]
		_, err := client.Admin().Windows().Updates().UpdatableAssets().ByUpdatableAssetId(groupId).Get(ctx, nil)
		return err
	})
}
