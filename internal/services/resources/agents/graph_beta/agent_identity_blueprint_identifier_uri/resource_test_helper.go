package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// AgentIdentityBlueprintIdentifierUriTestResource implements the types.TestResource interface
type AgentIdentityBlueprintIdentifierUriTestResource struct{}

// Exists checks whether the identifier URI exists on the application.
func (r AgentIdentityBlueprintIdentifierUriTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByStringArrayMembership(
		ctx,
		state,
		"blueprint_id",
		"identifierUris",
		"identifier_uri",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
			return client.Applications().ByApplicationId(parentID).Get(ctx, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "identifierUris"},
				},
			})
		},
	)
}
