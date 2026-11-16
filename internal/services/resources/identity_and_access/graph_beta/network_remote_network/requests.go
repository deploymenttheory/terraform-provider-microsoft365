package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const (
	collectionPath         = "/networkaccess/connectivity/remoteNetworks"
	itemPath               = collectionPath + "/{remoteNetworkId}"
	forwardingProfilesPath = itemPath + "/forwardingProfiles"
)

var (
	errMissingID            = errors.New("graph did not return an ID for the created remote network")
	errInvalidResponse      = errors.New("graph omitted required remote network fields")
	errNoForwardingProfiles = errors.New("at least one forwarding profile ID is required")
	errDeletionTimedOut     = errors.New("timed out waiting for remote network deletion to finish")
)

// deletionPollInterval is a variable so unit tests can shorten it.
var deletionPollInterval = 5 * time.Second

func profileRefs(ids []string) []map[string]string {
	refs := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		refs = append(refs, map[string]string{"id": id})
	}
	return refs
}

func (r *NetworkRemoteNetworkResource) create(ctx context.Context, name, region string, profileIDs []string) (*remoteNetworkResponse, error) {
	if len(profileIDs) == 0 {
		return nil, errNoForwardingProfiles
	}
	body := map[string]any{
		"name":               name,
		"region":             region,
		"forwardingProfiles": profileRefs(profileIDs),
	}
	var response remoteNetworkResponse
	if err := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.POST, collectionPath, nil, body, &response); err != nil {
		return nil, fmt.Errorf("create remote network: %w", err)
	}
	if response.ID == nil || *response.ID == "" {
		return nil, errMissingID
	}
	return &response, nil
}

// get uses the item GET because the forwardingProfiles navigation GET returns
// an empty collection even when profiles are associated.
func (r *NetworkRemoteNetworkResource) get(ctx context.Context, id string) (*remoteNetworkResponse, error) {
	var response remoteNetworkResponse
	if err := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.GET, itemPath, map[string]string{"remoteNetworkId": id}, nil, &response,
	); err != nil {
		return nil, fmt.Errorf("get remote network: %w", err)
	}
	return &response, nil
}

func (r *NetworkRemoteNetworkResource) updateName(ctx context.Context, id, name string) error {
	if err := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.PATCH, itemPath, map[string]string{"remoteNetworkId": id}, map[string]any{"name": name}, nil,
	); err != nil {
		return fmt.Errorf("update remote network name: %w", err)
	}
	return nil
}

// updateForwardingProfiles sends the full set as a $delta PATCH. The entity PATCH
// rejects forwardingProfiles, and the profile-side associations would overwrite
// other remote networks.
func (r *NetworkRemoteNetworkResource) updateForwardingProfiles(ctx context.Context, id string, profileIDs []string) error {
	if len(profileIDs) == 0 {
		return errNoForwardingProfiles
	}
	body := map[string]any{
		"@context": "#$delta",
		"value":    profileRefs(profileIDs),
	}
	if err := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.PATCH, forwardingProfilesPath, map[string]string{"remoteNetworkId": id}, body, nil,
	); err != nil {
		return fmt.Errorf("update remote network forwarding profiles: %w", err)
	}
	return nil
}

func (r *NetworkRemoteNetworkResource) delete(ctx context.Context, id string) error {
	if err := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.DELETE, itemPath, map[string]string{"remoteNetworkId": id}, nil, nil,
	); err != nil {
		return fmt.Errorf("delete remote network: %w", err)
	}
	return nil
}

// waitForDeletion waits for the asynchronous delete. Until it finishes, GET returns
// 404 "... is on deletion process." and the name cannot be reused (409).
func (r *NetworkRemoteNetworkResource) waitForDeletion(ctx context.Context, id string) error {
	for {
		_, err := r.get(ctx, id)
		if ctx.Err() != nil {
			return fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())
		}
		if err != nil {
			info := kiotaerrors.GraphError(ctx, err)
			switch info.StatusCode {
			case 404:
				if !strings.Contains(strings.ToLower(info.ErrorMessage), "deletion process") {
					return nil
				}
			case 401, 403:
				return err
			}
			// DELETE succeeded; keep polling through transient errors.
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())
		case <-time.After(deletionPollInterval):
		}
	}
}
