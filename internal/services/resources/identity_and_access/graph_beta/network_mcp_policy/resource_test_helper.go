package graphBetaNetworkMCPPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkMCPPolicyTestResource struct{}

func (r NetworkMCPPolicyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			resource := &NetworkMCPPolicyResource{client: client}
			readAdapter := resource.client.GetAdapter()
			readRequest, err := newMCPPolicyRequestInformation(
				ctx,
				readAdapter,
				abstractions.GET,
				state.ID,
				nil,
			)
			if err != nil {
				return err
			}
			readResult, err := readAdapter.SendPrimitive(ctx, readRequest, "[]byte", mcpPolicyErrorMapping)
			_, err = decodeSendMCPPolicyResponse(readResult, err)
			return err
		},
	)
}
