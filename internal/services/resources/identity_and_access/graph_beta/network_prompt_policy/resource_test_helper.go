package graphBetaNetworkPromptPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkPromptPolicyTestResource struct{}

func (r NetworkPromptPolicyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			resource := &NetworkPromptPolicyResource{client: client}
			readAdapter := resource.client.GetAdapter()
			readRequest, err := newPromptPolicyRequestInformation(
				ctx,
				readAdapter,
				abstractions.GET,
				state.ID,
				nil,
			)
			if err != nil {
				return err
			}
			readResult, err := readAdapter.Send(
				ctx,
				readRequest,
				createPromptPolicyResponseFromDiscriminatorValue,
				promptPolicyErrorMapping,
			)
			_, err = decodeSendPromptPolicyResponse(readResult, err)
			return err
		},
	)
}
