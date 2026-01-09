package graphBetaAgentsAgentCollection

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AgentCollectionTestResource implements the types.TestResource interface
type AgentCollectionTestResource struct{}

// Exists checks whether the agent collection exists in Microsoft Graph
func (r AgentCollectionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.AgentRegistry().AgentCollections().ByAgentCollectionId(state.ID).Get(ctx, nil)
		return err
	})
}
