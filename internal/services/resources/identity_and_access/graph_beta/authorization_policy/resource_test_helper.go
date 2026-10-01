package graphBetaAuthorizationPolicy

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

var errSingletonNotFound = errors.New("authorizationPolicy singleton missing from response")

// AuthorizationPolicyTestResource implements the types.TestResource interface for authorization policies.
type AuthorizationPolicyTestResource struct{}

// Exists checks whether the singleton authorization policy is readable in Microsoft Graph.
func (r AuthorizationPolicyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
			response, err := client.Policies().AuthorizationPolicy().Get(ctx, nil)
			if err != nil {
				return err
			}
			if response != nil {
				for _, policy := range response.GetValue() {
					if policy != nil && policy.GetId() != nil && *policy.GetId() == singletonID {
						return nil
					}
				}
			}
			return errSingletonNotFound
		},
	)
}
