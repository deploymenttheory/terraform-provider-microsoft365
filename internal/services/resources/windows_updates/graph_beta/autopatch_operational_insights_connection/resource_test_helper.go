package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WindowsUpdatesAutopatchOperationalInsightsConnectionTestResource struct{}

func (r WindowsUpdatesAutopatchOperationalInsightsConnectionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Admin().Windows().Updates().ResourceConnections().ByResourceConnectionId(state.ID).Get(ctx, nil)
		return err
	})
}
