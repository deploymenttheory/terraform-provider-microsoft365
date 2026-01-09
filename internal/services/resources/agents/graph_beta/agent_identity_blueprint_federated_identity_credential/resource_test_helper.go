package graphBetaAgentIdentityBlueprintFederatedIdentityCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AgentIdentityBlueprintFederatedIdentityCredentialTestResource implements the types.TestResource interface
type AgentIdentityBlueprintFederatedIdentityCredentialTestResource struct{}

// Exists checks whether the federated identity credential exists in Microsoft Graph
func (r AgentIdentityBlueprintFederatedIdentityCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"blueprint_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error {
			_, err := client.Applications().ByApplicationId(attributeValue).FederatedIdentityCredentials().ByFederatedIdentityCredentialId(resourceID).Get(ctx, nil)
			return err
		},
	)
}
