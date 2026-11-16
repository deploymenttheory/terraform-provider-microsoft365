package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// OnPremisesConnectorGroupMemberTestResource implements the types.TestResource interface for connector group memberships
type OnPremisesConnectorGroupMemberTestResource struct{}

// Exists checks whether the connector is still a member of the connector group in Terraform state
func (r OnPremisesConnectorGroupMemberTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		connectorID := state.Attributes["connector_id"]
		connectorGroupID := state.Attributes["connector_group_id"]

		memberOf, err := client.
			OnPremisesPublishingProfiles().
			ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
			Connectors().
			ByConnectorId(connectorID).
			MemberOf().
			Get(ctx, nil)
		if err != nil {
			return err
		}

		if memberOf != nil {
			for _, connectorGroup := range memberOf.GetValue() {
				if connectorGroup.GetId() != nil && strings.EqualFold(*connectorGroup.GetId(), connectorGroupID) {
					return nil
				}
			}
		}

		// Report a missing membership as a 404 so exists.CheckResourceExists
		// treats it as "does not exist" rather than as a failure.
		notFound := &abstractions.ApiError{Message: fmt.Sprintf("connector %s is not a member of connector group %s", connectorID, connectorGroupID)}
		notFound.SetStatusCode(404)
		return notFound
	})
}
