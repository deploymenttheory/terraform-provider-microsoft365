package graphBetaNetworkConditionalAccessSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

// NetworkConditionalAccessSettingsTestResource checks the singleton's availability.
type NetworkConditionalAccessSettingsTestResource struct{}

func (r NetworkConditionalAccessSettingsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // The generic helper adds operation context.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
		getConfig := customrequests.GetRequestConfig{Endpoint: "/networkAccess/settings/conditionalAccess", ResourceIDPattern: customrequests.ResourceIDPatternNone}
		_, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig)
		return err
	})
}
