package graphBetaRoleDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// RoleDefinitionTestResource implements the types.TestResource interface for role definitions
type RoleDefinitionTestResource struct{}

// Exists checks whether the role definition exists in Microsoft Graph
func (r RoleDefinitionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.RoleManagement().CloudPC().RoleDefinitions().ByUnifiedRoleDefinitionId(state.ID).Get(ctx, nil)
		return err
	})
}
