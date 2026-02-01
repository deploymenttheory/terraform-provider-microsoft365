package graphBetaApplicationFederatedIdentityCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ApplicationFederatedIdentityCredentialTestResource implements the types.TestResource interface
type ApplicationFederatedIdentityCredentialTestResource struct{}

// Exists checks whether the federated identity credential exists in Microsoft Graph
func (r ApplicationFederatedIdentityCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"application_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error {
			_, err := client.Applications().ByApplicationId(attributeValue).FederatedIdentityCredentials().ByFederatedIdentityCredentialId(resourceID).Get(ctx, nil)
			return err
		},
	)
}
