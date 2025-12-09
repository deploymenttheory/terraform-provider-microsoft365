package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// AddOwner adds an owner to the agent identity
// POST /servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post-owners?view=graph-rest-beta
func AddOwner(ctx context.Context, adapter abstractions.RequestAdapter, agentIdentityID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding owner %s to agent identity %s", ownerID, agentIdentityID))

	endpoint := fmt.Sprintf("servicePrincipals/%s/microsoft.graph.agentIdentity/owners/$ref", agentIdentityID)
	requestBody := ConstructAddOwnerOrSponsorRequest(ownerID)

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    endpoint,
		RequestBody: requestBody,
	}

	err := customrequests.PostRequestNoContent(ctx, adapter, config)
	if err != nil {
		return fmt.Errorf("failed to add owner %s: %w", ownerID, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully added owner %s", ownerID))
	return nil
}

// RemoveOwner removes an owner from the agent identity
// DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/{id}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-delete-owners?view=graph-rest-beta
func RemoveOwner(ctx context.Context, adapter abstractions.RequestAdapter, agentIdentityID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing owner %s from agent identity %s", ownerID, agentIdentityID))

	endpoint := fmt.Sprintf("servicePrincipals/%s/microsoft.graph.agentIdentity/owners", agentIdentityID)

	config := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          endpoint,
		ResourceID:        ownerID,
		ResourceIDPattern: "/{id}",
		EndpointSuffix:    "/$ref",
	}

	err := customrequests.DeleteRequestByResourceId(ctx, adapter, config)
	if err != nil {
		return fmt.Errorf("failed to remove owner %s: %w", ownerID, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed owner %s", ownerID))
	return nil
}
