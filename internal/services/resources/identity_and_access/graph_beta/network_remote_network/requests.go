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

// get reads the remote network. The item GET returns forwardingProfiles inline, whereas the
// forwardingProfiles navigation property GET returns an empty collection even when associated.
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

// updateForwardingProfiles replaces this remote network's profile associations with a $delta
// PATCH carrying the full desired set. Graph rejects forwardingProfiles on the entity PATCH, and
// writing the profile-side associations would clobber other remote networks.
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

// waitForDeletion polls until Graph has finished deleting the remote network.
// DELETE returns before the deletion completes (about 100s live). Until then,
// GET returns 404 "... is on deletion process." and the name cannot be reused,
// so a replacement that keeps the name would fail with 409.
func (r *NetworkRemoteNetworkResource) waitForDeletion(ctx context.Context, id string) error {
	for {
		_, err := r.get(ctx, id)
		if err != nil {
			info := kiotaerrors.GraphError(ctx, err)
			if info.StatusCode != 404 {
				return err
			}
			if !strings.Contains(strings.ToLower(info.ErrorMessage), "deletion process") {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())
		case <-time.After(deletionPollInterval):
		}
	}
}
