package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// AgentIdentityBlueprintCertificateCredentialTestResource implements the types.TestResource interface
type AgentIdentityBlueprintCertificateCredentialTestResource struct{}

// Exists checks whether the certificate credential exists in the application's keyCredentials
func (r AgentIdentityBlueprintCertificateCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByArrayMembership(
		ctx,
		state,
		"blueprint_id",
		"keyCredentials",
		"keyId",
		"key_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
			return client.Applications().ByApplicationId(parentID).Get(ctx, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "keyCredentials"},
				},
			})
		},
	)
}
