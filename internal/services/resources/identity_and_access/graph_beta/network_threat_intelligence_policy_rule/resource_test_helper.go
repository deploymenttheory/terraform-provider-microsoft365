package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkThreatIntelligencePolicyRuleTestResource struct{}

func (r NetworkThreatIntelligencePolicyRuleTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"threat_intelligence_policy_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parent, id string) error {
			resource := &NetworkThreatIntelligencePolicyRuleResource{client: client}
			_, err := resource.getThreatIntelligencePolicyRule(
				ctx,
				parent,
				id,
			)
			return err
		},
	)
}
