package graphBetaDirectoryRoleAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// DirectoryRoleAssignmentTestResource provides resource existence checks for acceptance tests.
type DirectoryRoleAssignmentTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r DirectoryRoleAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Preserve the shared acceptance helper's result.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(state.ID).Get(ctx, nil)
		return err
	})
}
