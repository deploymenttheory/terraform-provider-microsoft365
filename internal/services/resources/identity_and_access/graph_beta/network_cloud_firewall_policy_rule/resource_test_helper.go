package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// NetworkCloudFirewallPolicyRuleTestResource checks existence without treating HTTP 400 as deletion.
type NetworkCloudFirewallPolicyRuleTestResource struct{}

// Exists requires an actual HTTP 404 before reporting that cleanup succeeded.
func (r NetworkCloudFirewallPolicyRuleTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, fmt.Errorf("create acceptance Graph client: %w", err)
	}
	_, err = client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(state.Attributes["policy_id"]).PolicyRules().ByPolicyRuleId(state.ID).Get(ctx, nil)
	found := true
	if err != nil {
		if errors.GraphError(ctx, err).StatusCode != 404 {
			return nil, fmt.Errorf("check cloud firewall existence: %w", err)
		}
		found = false
	}
	return &found, nil
}
