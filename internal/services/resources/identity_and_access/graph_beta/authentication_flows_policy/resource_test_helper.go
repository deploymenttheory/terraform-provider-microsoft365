package graphBetaAuthenticationFlowsPolicy

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

var errSingletonNotFound = errors.New("authenticationFlowsPolicy singleton missing from response")

// AuthenticationFlowsPolicyTestResource implements the types.TestResource interface for authentication flows policies.
type AuthenticationFlowsPolicyTestResource struct{}

// Exists checks whether the singleton authentication flows policy is readable in Microsoft Graph.
func (r AuthenticationFlowsPolicyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
			response, err := client.Policies().AuthenticationFlowsPolicy().Get(ctx, nil)
			if err != nil {
				return err
			}
			if response != nil && response.GetId() != nil && *response.GetId() == singletonID {
				return nil
			}
			return errSingletonNotFound
		},
	)
}
