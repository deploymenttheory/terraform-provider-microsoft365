package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// AddSponsor adds a sponsor to the agent identity blueprint
// POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/sponsors/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-post-sponsors?view=graph-rest-beta
func AddSponsor(ctx context.Context, adapter abstractions.RequestAdapter, blueprintID, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding sponsor %s to blueprint %s", sponsorID, blueprintID))

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/sponsors/$ref", blueprintID)
	requestBody := ConstructAddOwnerOrSponsorRequest(sponsorID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: requestBody,
	}

	err := customrequests.PostRequestNoContent(ctx, adapter, config)
	if err != nil {
		return fmt.Errorf("failed to add sponsor %s: %w", sponsorID, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully added sponsor %s", sponsorID))
	return nil
}

// RemoveSponsor removes a sponsor from the agent identity blueprint
// DELETE /applications/{id}/microsoft.graph.agentIdentityBlueprint/sponsors/{sponsorObjectId}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-delete-sponsors?view=graph-rest-beta
func RemoveSponsor(ctx context.Context, adapter abstractions.RequestAdapter, blueprintID, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing sponsor %s from blueprint %s", sponsorID, blueprintID))

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/sponsors", blueprintID)

	config := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          endpoint,
		ResourceID:        sponsorID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/$ref",
	}

	err := customrequests.DeleteRequestByResourceId(ctx, adapter, config)
	if err != nil {
		return fmt.Errorf("failed to remove sponsor %s: %w", sponsorID, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed sponsor %s", sponsorID))
	return nil
}
