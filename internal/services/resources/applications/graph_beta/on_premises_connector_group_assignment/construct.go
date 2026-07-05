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

	// Microsoft Learn documents assignment as a reference update:
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta
	//
	// This intentionally differs from most relationship resources in the
	// provider, which often POST a concrete assignment object. Graph requires a
	// PUT to /applications/{application-id}/connectorGroup/$ref with only
	// @odata.id pointing at the connector group collection item.
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
	// Learn's example uses the Application Proxy connector group collection URL
	// as the reference target, not /applications/{id}/connectorGroup:
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta
	return fmt.Sprintf("https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/%s", connectorGroupID)
}

func compositeID(applicationID, connectorGroupID string) string {
	return fmt.Sprintf("%s/%s", applicationID, connectorGroupID)
}
