package graphBetaNetworkPromptPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkPromptPolicyRuleTestResource struct{}

func (r NetworkPromptPolicyRuleTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"prompt_policy_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, promptPolicyID, ruleID string) error {
			resource := &NetworkPromptPolicyRuleResource{client: client}
			_, err := resource.getPromptPolicyRule(
				ctx,
				promptPolicyID,
				ruleID,
			)
			return err
		},
	)
}
