package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

type OnPremisesConnectorGroupTestResource struct{}

func (r OnPremisesConnectorGroupTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		getConfig := customrequests.GetRequestConfig{
			APIVersion:        customrequests.GraphAPIBeta,
			Endpoint:          "onPremisesPublishingProfiles/applicationProxy/connectorGroups",
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID,
		}
		_, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig)
		return err
	})
}
