package graphBetaNetworkProxyAutoConfiguration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
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
			r.client = client
			var readResponse pacResponse
			readErr := customrequests.JSONRequest(
				ctx,
				r.client.GetAdapter(),
				abstractions.GET,
				r.ResourcePath+"/{pacId}",
				map[string]string{"pacId": state.ID},
				nil,
				&readResponse,
			)
			_, err := &readResponse, readErr
			if err != nil {
				err = fmt.Errorf("get custom PAC: %w", err)
			}
			return err
		},
	)
}
