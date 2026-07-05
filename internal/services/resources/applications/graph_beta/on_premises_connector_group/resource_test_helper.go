package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

type OnPremisesConnectorGroupTestResource struct{}

func buildConnectorGroupItemPath(connectorGroupID string) string {
	return fmt.Sprintf("onPremisesPublishingProfiles/applicationProxy/connectorGroups/%s", connectorGroupID)
}

func (r OnPremisesConnectorGroupTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		requestInfo := abstractions.NewRequestInformation()
		requestInfo.Method = abstractions.GET
		requestInfo.UrlTemplate = "{+baseurl}/" + buildConnectorGroupItemPath(state.ID)
		requestInfo.PathParameters = map[string]string{
			"baseurl": "https://graph.microsoft.com/beta",
		}

		errorMapping := abstractions.ErrorMappings{
			"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
		}

		_, err := client.GetAdapter().Send(
			ctx,
			requestInfo,
			graphmodels.CreateConnectorGroupFromDiscriminatorValue,
			errorMapping,
		)
		return err
	})
}
