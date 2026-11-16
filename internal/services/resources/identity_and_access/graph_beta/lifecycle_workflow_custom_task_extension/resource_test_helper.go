package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// LifecycleWorkflowCustomTaskExtensionTestResource implements the types.TestResource interface for custom task extensions
type LifecycleWorkflowCustomTaskExtensionTestResource struct{}

// Exists checks whether the custom task extension exists in Microsoft Graph
func (r LifecycleWorkflowCustomTaskExtensionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.IdentityGovernance().LifecycleWorkflows().CustomTaskExtensions().ByCustomTaskExtensionId(state.ID).Get(ctx, nil)
		return err
	})
}
