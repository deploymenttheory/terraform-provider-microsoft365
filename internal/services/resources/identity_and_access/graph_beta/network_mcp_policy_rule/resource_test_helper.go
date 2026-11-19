package graphBetaNetworkMCPPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkMCPPolicyRuleTestResource struct{}

func (r NetworkMCPPolicyRuleTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			resource := &NetworkMCPPolicyRuleResource{client: client}
			readAdapter := resource.client.GetAdapter()
			readRequest, err := newMCPPolicyRuleRequestInformation(
				ctx,
				readAdapter,
				abstractions.GET,
				state.Attributes["mcp_policy_id"],
				state.ID,
				nil,
			)
			if err != nil {
				return err
			}
			readResult, err := readAdapter.SendPrimitive(ctx, readRequest, "[]byte", mcpPolicyRuleErrorMapping)
			_, err = decodeSendMCPPolicyRuleResponse(readResult, err)
			return err
		},
	)
}
