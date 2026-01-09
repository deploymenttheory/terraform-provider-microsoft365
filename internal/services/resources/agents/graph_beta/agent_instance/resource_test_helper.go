package graphBetaAgentInstance

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// AgentInstanceTestResource implements the types.TestResource interface
type AgentInstanceTestResource struct{}

// Exists checks whether the agent instance exists in Microsoft Graph
func (r AgentInstanceTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.AgentRegistry().AgentInstances().ByAgentInstanceId(state.ID).Get(ctx, nil)
		return err
	})
}
