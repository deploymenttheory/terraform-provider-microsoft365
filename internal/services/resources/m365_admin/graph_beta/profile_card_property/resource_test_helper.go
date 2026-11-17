package graphBetaM365AdminProfileCardProperty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// ProfileCardPropertyTestResource implements the types.TestResource interface for profile card properties
type ProfileCardPropertyTestResource struct{}

// Exists checks whether the profile card property exists in Microsoft Graph
func (r ProfileCardPropertyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Admin().People().ProfileCardProperties().ByProfileCardPropertyId(state.ID).Get(ctx, nil)
		return err
	})
}
