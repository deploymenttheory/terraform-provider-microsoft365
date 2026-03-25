package graphBetaAdministrativeUnitMembership

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type AdministrativeUnitMembershipTestResource struct{}

func (r AdministrativeUnitMembershipTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		administrativeUnitID := state.Attributes["administrative_unit_id"]
		_, err := client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Get(ctx, nil)
		return err
	})
}
