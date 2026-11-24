package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

type NetworkPromptPolicyRuleTestResource struct{}

func (r NetworkPromptPolicyRuleTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // The existence helper treats the raw Graph error as not found or failed.
	getRule := func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, promptPolicyID, ruleID string) error {
		getConfig := customrequests.GetRequestConfig{Endpoint: fmt.Sprintf(promptPolicyRuleEndpoint, promptPolicyID, ruleID)}
		_, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig)
		return err
	}
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExistsByCompositeID(ctx, state, "prompt_policy_id", getRule)
}
