package graphBetaGroupAppRoleAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GroupAppRoleAssignmentTestResource implements the types.TestResource interface for group app role assignments
type GroupAppRoleAssignmentTestResource struct{}

// Exists checks whether the group app role assignment exists in Microsoft Graph
func (r GroupAppRoleAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		groupID := state.Attributes["target_group_id"]
		_, err := client.
			Groups().
			ByGroupId(groupID).
			AppRoleAssignments().
			ByAppRoleAssignmentId(state.ID).
			Get(ctx, nil)
		return err
	})
}
