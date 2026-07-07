package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type OnPremisesConnectorGroupAssignmentTestResource struct{}

func (r OnPremisesConnectorGroupAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		applicationID := state.Attributes["application_id"]
		connectorGroupID := state.Attributes["connector_group_id"]

		connectorGroup, err := client.Applications().ByApplicationId(applicationID).ConnectorGroup().Get(ctx, nil)
		if err != nil {
			return err
		}
		if connectorGroup == nil || connectorGroup.GetId() == nil || *connectorGroup.GetId() != connectorGroupID {
			return fmt.Errorf("connector group %s is not assigned to application %s", connectorGroupID, applicationID)
		}
		return nil
	})
}
