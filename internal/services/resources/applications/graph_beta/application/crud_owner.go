package graphBetaApplication

import (
	"context"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAddOwnerRequest constructs the request body for adding an owner
// The request body must contain an @odata.id property with the URL of the user or service principal
// REF: https://learn.microsoft.com/en-us/graph/api/application-post-owners?view=graph-rest-beta
func ConstructAddOwnerRequest(userID string) graphmodels.ReferenceCreateable {
	requestBody := graphmodels.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/beta/directoryObjects/%s", userID)
	requestBody.SetOdataId(&odataId)
	return requestBody
}

// AddOwner adds an owner to the application
// POST /applications/{id}/owners/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/application-post-owners?view=graph-rest-beta
func AddOwner(ctx context.Context, adapter abstractions.RequestAdapter, applicationID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Adding owner %s to application %s", ownerID, applicationID))

	endpoint := fmt.Sprintf("applications/%s/owners/$ref", applicationID)
	requestBody := ConstructAddOwnerRequest(ownerID)

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

// RemoveOwner removes an owner from the application
// DELETE /applications/{id}/owners/{ownerId}/$ref
// REF: https://learn.microsoft.com/en-us/graph/api/application-delete-owners?view=graph-rest-beta
func RemoveOwner(ctx context.Context, adapter abstractions.RequestAdapter, applicationID, ownerID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Removing owner %s from application %s", ownerID, applicationID))

	endpoint := fmt.Sprintf("applications/%s/owners", applicationID)

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
