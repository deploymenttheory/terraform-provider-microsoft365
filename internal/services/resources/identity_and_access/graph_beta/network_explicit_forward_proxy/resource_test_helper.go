package graphBetaNetworkExplicitForwardProxy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

type NetworkExplicitForwardProxyTestResource struct{}

func (r NetworkExplicitForwardProxyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared existence helper supplies operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
			r := NewNetworkExplicitForwardProxyResource().(*NetworkExplicitForwardProxyResource)
			getConfig := customrequests.GetRequestConfig{
				Endpoint:          r.ResourcePath,
				ResourceIDPattern: customrequests.ResourceIDPatternNone,
			}
			if _, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig); err != nil {
				return fmt.Errorf("get explicit proxy: %w", err)
			}
			return nil
		},
	)
}
