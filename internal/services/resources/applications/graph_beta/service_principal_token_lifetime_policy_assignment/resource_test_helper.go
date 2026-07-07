package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ServicePrincipalTokenLifetimePolicyAssignmentTestResource implements the types.TestResource interface
type ServicePrincipalTokenLifetimePolicyAssignmentTestResource struct{}

// Exists checks whether the token lifetime policy is assigned to the service principal in Microsoft Graph
func (r ServicePrincipalTokenLifetimePolicyAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		// ID format: service_principal_id/token_lifetime_policy_id
		spID := state.Attributes["service_principal_id"]
		policyID := state.Attributes["token_lifetime_policy_id"]

		if spID == "" || policyID == "" {
			return fmt.Errorf("missing service_principal_id or token_lifetime_policy_id in state")
		}

		policies, err := client.ServicePrincipals().ByServicePrincipalId(spID).TokenLifetimePolicies().Get(ctx, nil)
		if err != nil {
			return err
		}

		for _, policy := range policies.GetValue() {
			if policy.GetId() != nil && *policy.GetId() == policyID {
				return nil
			}
		}

		return fmt.Errorf("token lifetime policy %s not assigned to service principal %s", policyID, spID)
	})
}
