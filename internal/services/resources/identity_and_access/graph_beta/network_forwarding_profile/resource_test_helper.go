package graphBetaNetworkForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// NetworkForwardingProfileTestResource provides resource existence checks for acceptance tests.
type NetworkForwardingProfileTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r NetworkForwardingProfileTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared existence helper supplies operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			r := NewNetworkForwardingProfileResource().(*NetworkForwardingProfileResource)
			r.client = client
			_, err := r.get(ctx, state.ID)
			return err
		},
	)
}
