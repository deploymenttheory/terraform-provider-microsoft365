package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// LifecycleWorkflowTestResource implements the types.TestResource interface for lifecycle workflows
type LifecycleWorkflowTestResource struct{}

// Exists checks whether the lifecycle workflow exists in Microsoft Graph
func (r LifecycleWorkflowTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.IdentityGovernance().LifecycleWorkflows().Workflows().ByWorkflowId(state.ID).Get(ctx, nil)
		return err
	})
}
