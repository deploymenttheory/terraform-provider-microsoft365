package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

const defaultGraphBetaBaseURL = "https://graph.microsoft.com/beta"

// errConnectorGroupIDRequired is returned when no connector group ID is supplied.
var errConnectorGroupIDRequired = errors.New("connector_group_id is required")

// constructResource builds the $ref body that assigns a connector to a connector group.
//
// Microsoft Learn documents the request as POST
// /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref
// with only @odata.id pointing at the connector group:
// https://learn.microsoft.com/en-us/graph/api/connector-post-memberof?view=graph-rest-beta
func constructResource(ctx context.Context, connectorGroupID string, graphBetaBaseURL string) (graphmodels.ReferenceCreateable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if connectorGroupID == "" {
		return nil, errConnectorGroupIDRequired
	}

	requestBody := graphmodels.NewReferenceCreate()
	odataID := connectorGroupODataID(graphBetaBaseURL, connectorGroupID)
	requestBody.SetOdataId(&odataID)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// connectorGroupODataID returns the connector group reference URL. It uses the
// client's configured base URL so national clouds reference their own endpoint.
func connectorGroupODataID(graphBetaBaseURL, connectorGroupID string) string {
	baseURL := strings.TrimRight(graphBetaBaseURL, "/")
	if baseURL == "" {
		baseURL = defaultGraphBetaBaseURL
	}

	return fmt.Sprintf("%s/onPremisesPublishingProfiles/applicationProxy/connectorGroups/%s", baseURL, connectorGroupID)
}
