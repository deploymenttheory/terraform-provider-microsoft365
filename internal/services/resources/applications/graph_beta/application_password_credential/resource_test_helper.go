package graphBetaApplicationPasswordCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// ApplicationPasswordCredentialTestResource implements the types.TestResource interface
type ApplicationPasswordCredentialTestResource struct{}

// Exists checks whether the password credential exists in Microsoft Graph.
// Note: Password credentials cannot be directly queried by keyId, so we list all
// credentials on the application and check if the keyId exists.
func (r ApplicationPasswordCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByArrayMembership(
		ctx,
		state,
		"application_id",
		"passwordCredentials",
		"keyId",
		"key_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
			return client.Applications().ByApplicationId(parentID).Get(ctx, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "passwordCredentials"},
				},
			})
		},
	)
}
