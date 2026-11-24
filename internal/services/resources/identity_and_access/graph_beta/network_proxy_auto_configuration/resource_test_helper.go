package graphBetaNetworkProxyAutoConfiguration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
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
			getConfig := customrequests.GetRequestConfig{
				Endpoint:          r.ResourcePath,
				ResourceIDPattern: customrequests.ResourceIDPatternSegment,
				ResourceID:        state.ID,
			}
			if _, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig); err != nil {
				return fmt.Errorf("get custom PAC %s: %w", state.ID, err)
			}
			return nil
		},
	)
}
