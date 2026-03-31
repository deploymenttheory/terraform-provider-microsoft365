package graphBetaWindowsUpdatesAutopatchRing

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WindowsUpdatesAutopatchRingTestResource struct{}

func (r WindowsUpdatesAutopatchRingTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		policyId := state.Attributes["policy_id"]
		ringId := state.ID
		_, err := client.Admin().Windows().Updates().Policies().ByPolicyId(policyId).Rings().ByRingId(ringId).Get(ctx, nil)
		return err
	})
}
