package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *OnPremisesConnectorGroupAssignmentResourceModel) (graphmodels.ReferenceUpdateable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewReferenceUpdate()
	odataID := connectorGroupODataID(data.ConnectorGroupID.ValueString())
	requestBody.SetOdataId(&odataID)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func connectorGroupODataID(connectorGroupID string) string {
	return fmt.Sprintf("https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/%s", connectorGroupID)
}

func compositeID(applicationID, connectorGroupID string) string {
	return fmt.Sprintf("%s/%s", applicationID, connectorGroupID)
}
