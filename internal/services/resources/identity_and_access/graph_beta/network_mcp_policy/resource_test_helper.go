package graphBetaNetworkMCPPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

type NetworkMCPPolicyTestResource struct{}

func (r NetworkMCPPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		getConfig := customrequests.GetRequestConfig{Endpoint: mcpPoliciesEndpoint + "/" + state.ID}
		_, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig)
		return err
	})
}
