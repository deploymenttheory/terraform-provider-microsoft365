package graphConditionalAccessTermsOfUse

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ConditionalAccessTermsOfUseTestResource implements the types.TestResource interface for terms of use
type ConditionalAccessTermsOfUseTestResource struct{}

// Exists checks whether the conditional access terms of use exists in Microsoft Graph
func (r ConditionalAccessTermsOfUseTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Agreements().ByAgreementId(state.ID).Get(ctx, nil)
		return err
	})
}
