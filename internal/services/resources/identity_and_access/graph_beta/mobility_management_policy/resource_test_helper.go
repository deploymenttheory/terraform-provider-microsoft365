package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// MobilityManagementPolicyTestResource implements the types.TestResource interface for mobility management policies
type MobilityManagementPolicyTestResource struct{}

// Exists checks whether the mobility management policy exists in Microsoft Graph
func (r MobilityManagementPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		policyType, policyID, err := parseResourceID(state.ID)
		if err != nil {
			return err
		}
		if policyType == policyTypeApp {
			_, err = client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).Get(ctx, nil)
			return err
		}
		_, err = client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).Get(ctx, nil)
		return err
	})
}
