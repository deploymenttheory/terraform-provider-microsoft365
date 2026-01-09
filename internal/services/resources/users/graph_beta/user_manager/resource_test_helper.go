package graphBetaUsersUserManager

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// UserManagerTestResource implements the types.TestResource interface for user manager relationships
type UserManagerTestResource struct{}

// Exists checks whether the user manager relationship exists in Microsoft Graph
func (r UserManagerTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.
			Users().
			ByUserId(state.ID).
			Manager().
			Get(ctx, nil)
		return err
	})
}
