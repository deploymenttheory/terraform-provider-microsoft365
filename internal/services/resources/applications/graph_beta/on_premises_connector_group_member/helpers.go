package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"fmt"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// onPremisesPublishingProfileID is the publishing profile that owns private
// network (Application Proxy) connectors and connector groups.
const onPremisesPublishingProfileID = "applicationProxy"

// listConnectorMemberships returns the connector groups the connector belongs to.
// A connector belongs to exactly one connector group, so callers use the first.
//
// API Calls:
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-list-memberof?view=graph-rest-beta
func (r *OnPremisesConnectorGroupMemberResource) listConnectorMemberships(ctx context.Context, connectorID string) ([]graphmodels.ConnectorGroupable, error) {
	memberOf, err := r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		Connectors().
		ByConnectorId(connectorID).
		MemberOf().
		Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("listing connector groups for connector %s: %w", connectorID, err)
	}

	if memberOf == nil {
		return nil, nil
	}

	return memberOf.GetValue(), nil
}

// listConnectorGroups returns every connector group in the tenant.
//
// API Calls:
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectorGroups
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connectorgroup-list?view=graph-rest-beta
func (r *OnPremisesConnectorGroupMemberResource) listConnectorGroups(ctx context.Context) ([]graphmodels.ConnectorGroupable, error) {
	connectorGroups, err := r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		ConnectorGroups().
		Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("listing connector groups: %w", err)
	}

	if connectorGroups == nil {
		return nil, nil
	}

	return connectorGroups.GetValue(), nil
}

// firstConnectorGroup returns the first connector group, or nil when there is none.
func firstConnectorGroup(connectorGroups []graphmodels.ConnectorGroupable) graphmodels.ConnectorGroupable {
	if len(connectorGroups) == 0 {
		return nil
	}

	return connectorGroups[0]
}

// selectDefaultConnectorGroup picks the tenant's default connector group.
func selectDefaultConnectorGroup(connectorGroups []graphmodels.ConnectorGroupable) graphmodels.ConnectorGroupable {
	for _, connectorGroup := range connectorGroups {
		if connectorGroup != nil && connectorGroup.GetIsDefault() != nil && *connectorGroup.GetIsDefault() {
			return connectorGroup
		}
	}

	return nil
}

// connectorGroupIDOf returns the connector group's ID, or an empty string.
func connectorGroupIDOf(connectorGroup graphmodels.ConnectorGroupable) string {
	if connectorGroup == nil || connectorGroup.GetId() == nil {
		return ""
	}

	return *connectorGroup.GetId()
}

// sameID compares Graph object IDs case-insensitively.
func sameID(a, b string) bool {
	return a != "" && strings.EqualFold(a, b)
}
