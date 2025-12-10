package graphBetaAgentUser

import (
	"context"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAddSponsorRequest constructs the request body for adding a sponsor
// The request body must contain an @odata.id property with the URL of the user or group
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-sponsors?view=graph-rest-beta
func ConstructAddSponsorRequest(sponsorID string) graphmodels.ReferenceCreateable {
	requestBody := graphmodels.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", sponsorID)
	requestBody.SetOdataId(&odataId)
	return requestBody
}

// AddSponsor adds a sponsor to the agent user
// POST /users/{usersId}/sponsors/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-sponsors?view=graph-rest-beta
func AddSponsor(ctx context.Context, adapter abstractions.RequestAdapter, agentUserId, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding sponsor %s to agent user %s", sponsorID, agentUserId))

	endpoint := fmt.Sprintf("/users/%s/sponsors/$ref", agentUserId)
	requestBody := ConstructAddSponsorRequest(sponsorID)

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

// RemoveSponsor removes a sponsor from the agent user
// DELETE /users/{usersId}/sponsors/{id}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-delete-sponsors?view=graph-rest-beta
func RemoveSponsor(ctx context.Context, adapter abstractions.RequestAdapter, agentUserId, sponsorID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing sponsor %s from agent user %s", sponsorID, agentUserId))

	endpoint := fmt.Sprintf("/users/%s/sponsors", agentUserId)

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
