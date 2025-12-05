package graphBetaApplicationsAgentIdentityBlueprint

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
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-post-owners?view=graph-rest-beta
func ConstructAddOwnerOrSponsorRequest(userID string) graphmodels.ReferenceCreateable {
	requestBody := graphmodels.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", userID)
	requestBody.SetOdataId(&odataId)
	return requestBody
}

// AddOwner adds an owner to the agent identity blueprint
// POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/owners/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-post-owners?view=graph-rest-beta
func AddOwner(ctx context.Context, adapter abstractions.RequestAdapter, blueprintID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding owner %s to blueprint %s", ownerID, blueprintID))

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/owners/$ref", blueprintID)
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

// RemoveOwner removes an owner from the agent identity blueprint
// DELETE /applications/{id}/microsoft.graph.agentIdentityBlueprint/owners/{id}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-delete-owners?view=graph-rest-beta
func RemoveOwner(ctx context.Context, adapter abstractions.RequestAdapter, blueprintID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing owner %s from blueprint %s", ownerID, blueprintID))

	endpoint := fmt.Sprintf("applications/%s/microsoft.graph.agentIdentityBlueprint/owners", blueprintID)

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
