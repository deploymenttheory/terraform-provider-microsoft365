package graphBetaAgentIdentity

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AgentIdentityTestResource implements the types.TestResource interface
type AgentIdentityTestResource struct{}

// Exists checks whether the agent identity exists in Microsoft Graph
func (r AgentIdentityTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.ServicePrincipals().ByServicePrincipalId(state.ID).Get(ctx, nil)
		return err
	})
}
