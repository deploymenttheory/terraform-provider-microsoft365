package graphBetaApplication

import (
	"context"
	"fmt"
	"strings"
	"time"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
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
	// A directory replica can briefly miss this application while its owner
	// object is being deleted. A failed DELETE has no reference side effect.
	for attempt := 0; err != nil && attempt < 15; attempt++ {
		info := errors.GraphError(ctx, err)
		if info.StatusCode != 404 || !strings.Contains(info.ErrorMessage, applicationID) {
			break
		}
		absent, verificationErr := ownerReferenceAbsent(ctx, adapter, applicationID, ownerID)
		if verificationErr == nil && absent {
			return nil
		}
		if verificationErr != nil && errors.GraphError(ctx, verificationErr).StatusCode != 404 {
			return fmt.Errorf("failed to verify removal of owner %s: %w", ownerID, verificationErr)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting to remove owner %s: %w", ownerID, ctx.Err())
		case <-time.After(2 * time.Second):
		}
		err = customrequests.DeleteRequestByResourceId(ctx, adapter, config)
	}
	if err != nil {
		info := errors.GraphError(ctx, err)
		// Deleting the owner directory object also removes its reference. Graph
		// returns this specific 400 when Terraform subsequently removes the link.
		if info.StatusCode == 400 && strings.Contains(info.ErrorMessage, "removed object references do not exist") && strings.Contains(info.ErrorMessage, "'owners'") {
			return nil
		}
		return fmt.Errorf("failed to remove owner %s: %w", ownerID, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed owner %s", ownerID))
	return nil
}

// ownerReferenceAbsent confirms the target is absent from every owners page.
// A DELETE 404 alone cannot distinguish a missing link from a missing replica.
func ownerReferenceAbsent(ctx context.Context, adapter abstractions.RequestAdapter, applicationID, ownerID string) (bool, error) {
	builder := msgraphbetasdk.NewGraphServiceClient(adapter).Applications().ByApplicationId(applicationID).Owners()
	for {
		page, err := builder.Get(ctx, nil)
		if err != nil {
			return false, fmt.Errorf("read application owners: %w", err)
		}
		if page == nil {
			return false, nil
		}
		for _, owner := range page.GetValue() {
			if owner != nil && owner.GetId() != nil && *owner.GetId() == ownerID {
				return false, nil
			}
		}
		next := page.GetOdataNextLink()
		if next == nil || *next == "" {
			return true, nil
		}
		builder = builder.WithUrl(*next)
	}
}
