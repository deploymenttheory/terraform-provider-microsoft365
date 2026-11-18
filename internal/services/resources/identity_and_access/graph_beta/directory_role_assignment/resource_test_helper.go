package graphBetaDirectoryRoleAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type DirectoryRoleAssignmentTestResource struct{}

func (r DirectoryRoleAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Preserve the shared acceptance helper's result.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(state.ID).Get(ctx, nil)
		return err
	})
}
