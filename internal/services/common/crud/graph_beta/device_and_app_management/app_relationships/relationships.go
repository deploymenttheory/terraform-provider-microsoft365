// Package apprelationships handles the relationship collection shared by Intune apps.
package apprelationships

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var (
	ErrIncompleteCollection = errors.New("graph returned an incomplete app relationship collection")
	ErrNotFound             = errors.New("app relationship not found (404)")
	appLocks                sync.Map
)

// Lock serializes collection writes across dependency and supersedence resources.
// Callers hold the lock through read-back so the next writer sees the completed write.
func Lock(ctx context.Context, sourceID string) (func(), error) {
	value, _ := appLocks.LoadOrStore(strings.ToLower(sourceID), make(chan struct{}, 1))
	lock := value.(chan struct{})
	select {
	case lock <- struct{}{}:
		return func() { <-lock }, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for app relationship lock: %w", ctx.Err())
	}
}

// SourceID extracts the app ID from Graph's source_target relationship identifier.
// Opaque IDs are deliberately left intact for the legacy ID-only lookup.
func SourceID(relationshipID string) string {
	source, target, ok := strings.Cut(relationshipID, "_")
	if !ok {
		return ""
	}
	if _, err := uuid.Parse(source); err != nil {
		return ""
	}
	if _, err := uuid.Parse(target); err != nil {
		return ""
	}
	return source
}

// List retrieves all pages before callers make any decisions about absence or replacement.
func List(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	sourceID string,
) ([]graphmodels.MobileAppRelationshipable, error) {
	builder := client.DeviceAppManagement().MobileApps().ByMobileAppId(sourceID).Relationships()
	page, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list app relationships: %w", err)
	}
	relationships := make([]graphmodels.MobileAppRelationshipable, 0)
	for {
		if page == nil || page.GetValue() == nil {
			return nil, ErrIncompleteCollection
		}
		for _, relationship := range page.GetValue() {
			if relationship == nil {
				return nil, ErrIncompleteCollection
			}
			relationships = append(relationships, relationship)
		}
		nextLink := page.GetOdataNextLink()
		if nextLink == nil || *nextLink == "" {
			return relationships, nil
		}
		page, err = builder.WithUrl(*nextLink).Get(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("read next app relationship page: %w", err)
		}
	}
}

// Incoming identifies relationships owned by another app.
func Incoming(relationship graphmodels.MobileAppRelationshipable) bool {
	return relationship.GetTargetType() != nil &&
		*relationship.GetTargetType() == graphmodels.PARENT_MOBILEAPPRELATIONSHIPTYPE
}

// Find matches a saved ID, or discovers the ID by target after a collection update.
func Find(
	relationships []graphmodels.MobileAppRelationshipable,
	relationshipID, targetID string,
) (graphmodels.MobileAppRelationshipable, error) {
	for _, relationship := range relationships {
		if Incoming(relationship) {
			continue
		}
		if relationshipID != "" {
			if relationship.GetId() != nil &&
				strings.EqualFold(*relationship.GetId(), relationshipID) {
				return relationship, nil
			}
		} else if relationship.GetTargetId() != nil && strings.EqualFold(*relationship.GetTargetId(), targetID) {
			return relationship, nil
		}
	}
	return nil, ErrNotFound
}
