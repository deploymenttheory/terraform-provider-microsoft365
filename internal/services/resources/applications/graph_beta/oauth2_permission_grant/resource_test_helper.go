package graphBetaOAuth2PermissionGrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// OAuth2PermissionGrantTestResource provides resource existence checks for acceptance tests.
type OAuth2PermissionGrantTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r OAuth2PermissionGrantTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Pass through the generic acceptance helper.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(state.ID).Get(ctx, nil)
		return err
	})
}
