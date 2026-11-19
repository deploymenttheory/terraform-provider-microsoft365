package graphBetaNetworkDeviceForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// NetworkDeviceForwardingProfileTestResource provides resource existence checks for acceptance tests.
type NetworkDeviceForwardingProfileTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r NetworkDeviceForwardingProfileTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared existence helper supplies operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			r := NewNetworkDeviceForwardingProfileResource().(*NetworkDeviceForwardingProfileResource)
			r.client = client
			_, err := r.get(ctx, state.ID)
			return err
		},
	)
}
