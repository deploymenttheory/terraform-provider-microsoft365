package graphBetaNetworkConditionalAccessSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// NetworkConditionalAccessSettingsTestResource checks the singleton's availability.
type NetworkConditionalAccessSettingsTestResource struct{}

func (r NetworkConditionalAccessSettingsTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic helper adds operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, _ *terraform.InstanceState) error {
			r := &NetworkConditionalAccessSettingsResource{client: client}
			request, err := r.client.NetworkAccess().Settings().ConditionalAccess().ToGetRequestInformation(ctx, nil)
			if err != nil {
				return err
			}
			result, err := r.client.GetAdapter().SendPrimitive(ctx, request, "[]byte", abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
			_, err = decodeSettingsResponse(result, err)
			return err
		},
	)
}
