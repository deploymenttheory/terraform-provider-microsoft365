package graphBetaNetworkProxyAutoConfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkProxyAutoConfigurationTestResource struct{}

func (r NetworkProxyAutoConfigurationTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared existence helper supplies operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			r := NewNetworkProxyAutoConfigurationResource().(*NetworkProxyAutoConfigurationResource)
			r.client = client
			_, err := r.get(ctx, state.ID)
			return err
		},
	)
}
