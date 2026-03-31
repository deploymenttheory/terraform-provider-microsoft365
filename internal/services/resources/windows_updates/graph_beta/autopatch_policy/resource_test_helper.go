package graphBetaWindowsUpdatesAutopatchPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// WindowsAutopatchPolicyTestResource implements the types.TestResource interface for Windows Autopatch policies
type WindowsAutopatchPolicyTestResource struct{}

// Exists checks whether the Windows Autopatch policy exists in Microsoft Graph
func (r WindowsAutopatchPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.Admin().Windows().Updates().Policies().ByPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
