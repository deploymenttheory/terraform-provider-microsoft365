package graphBetaAdministrativeUnitRoleAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type AdministrativeUnitRoleAssignmentTestResource struct{}

func (r AdministrativeUnitRoleAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		administrativeUnitID := state.Attributes["administrative_unit_id"]
		scopedRoleMembershipID := state.Attributes["id"]
		_, err := client.
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			ScopedRoleMembers().
			ByScopedRoleMembershipId(scopedRoleMembershipID).
			Get(ctx, nil)
		return err
	})
}
