package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// B2bManagementPolicyTestResource implements the types.TestResource interface for B2B management policies
type B2bManagementPolicyTestResource struct{}

// Exists checks whether the B2B management policy exists in Microsoft Graph
func (r B2bManagementPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, b2bManagementPolicyExists)
}

func b2bManagementPolicyExists(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
	getPolicy := func(ctx context.Context) (graphmodels.B2bManagementPolicyable, error) {
		policy, err := client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(state.ID).Get(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("checking B2B management policy %s: %w", state.ID, err)
		}
		return policy, nil
	}
	_, err := getPolicy(ctx)
	if err != nil && errors.GraphError(ctx, err).StatusCode == 404 {
		_, err = confirmPolicyNotFound(ctx, state.ID, err, getPolicy)
	}
	return err
}
