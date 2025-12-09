package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAddOwnerOrSponsorRequest constructs the request body for adding an owner or sponsor
// The request body must contain an @odata.id property with the URL of the user or service principal
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post-sponsors?view=graph-rest-beta
func ConstructAddOwnerOrSponsorRequest(userID string) graphmodels.ReferenceCreateable {
	requestBody := graphmodels.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", userID)
	requestBody.SetOdataId(&odataId)
	return requestBody
}

// AddSponsor adds a sponsor to the agent identity
// POST /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-post-sponsors?view=graph-rest-beta
func AddSponsor(ctx context.Context, adapter abstractions.RequestAdapter, agentIdentityID, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding sponsor %s to agent identity %s", sponsorID, agentIdentityID))

	endpoint := fmt.Sprintf("servicePrincipals/%s/microsoft.graph.agentIdentity/sponsors/$ref", agentIdentityID)
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

// RemoveSponsor removes a sponsor from the agent identity
// DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/{sponsorObjectId}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentity-delete-sponsors?view=graph-rest-beta
func RemoveSponsor(ctx context.Context, adapter abstractions.RequestAdapter, agentIdentityID, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing sponsor %s from agent identity %s", sponsorID, agentIdentityID))

	endpoint := fmt.Sprintf("servicePrincipals/%s/microsoft.graph.agentIdentity/sponsors", agentIdentityID)

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
