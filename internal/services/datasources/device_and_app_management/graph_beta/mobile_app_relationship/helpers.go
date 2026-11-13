package graphBetaMobileAppRelationship

import (
	"errors"
	"fmt"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var (
	errRelationshipNotFound  = errors.New("no mobile app relationship found")
	errMultipleRelationships = errors.New(
		"the lookup must match exactly one mobile app relationship",
	)
)

// lookupMethod represents the different ways to look up a relationship.
type lookupMethod int

const (
	lookupByRelationshipId lookupMethod = iota
	lookupByTargetId
	lookupByTargetDisplayName
	lookupByODataQuery
)

// determineLookupMethod selects the lookup from the configured attribute; schema
// validators guarantee exactly one is set.
func determineLookupMethod(object MobileAppRelationshipDataSourceModel) lookupMethod {
	switch {
	case object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case object.TargetID.ValueString() != "":
		return lookupByTargetId
	case object.TargetDisplayName.ValueString() != "":
		return lookupByTargetDisplayName
	default:
		return lookupByRelationshipId
	}
}

// relationshipMatcher reports whether a relationship satisfies the configured lookup.
type relationshipMatcher func(graphmodels.MobileAppRelationshipable) bool

func matchAll(graphmodels.MobileAppRelationshipable) bool { return true }

func matchId(id string) relationshipMatcher {
	return func(r graphmodels.MobileAppRelationshipable) bool {
		return r.GetId() != nil && strings.EqualFold(*r.GetId(), id)
	}
}

func matchTargetId(targetID string) relationshipMatcher {
	return func(r graphmodels.MobileAppRelationshipable) bool {
		return r.GetTargetId() != nil && strings.EqualFold(*r.GetTargetId(), targetID)
	}
}

func matchTargetDisplayName(name string) relationshipMatcher {
	return func(r graphmodels.MobileAppRelationshipable) bool {
		return r.GetTargetDisplayName() != nil && strings.EqualFold(*r.GetTargetDisplayName(), name)
	}
}

// selectRelationship applies the lookup and optional target_type filter, and
// requires exactly one match.
func selectRelationship(
	relationships []graphmodels.MobileAppRelationshipable,
	match relationshipMatcher,
	targetType, criteria string,
) (graphmodels.MobileAppRelationshipable, error) {
	var matches []graphmodels.MobileAppRelationshipable
	for _, relationship := range relationships {
		if relationship == nil || !match(relationship) {
			continue
		}
		if targetType != "" && (relationship.GetTargetType() == nil ||
			relationship.GetTargetType().String() != targetType) {
			continue
		}
		matches = append(matches, relationship)
	}

	if targetType != "" {
		criteria = fmt.Sprintf("%s and target_type %q", criteria, targetType)
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("%w with %s", errRelationshipNotFound, criteria)
	case 1:
		return matches[0], nil
	default:
		return nil, fmt.Errorf(
			"%w: found %d with %s; narrow it with target_type or a more specific lookup",
			errMultipleRelationships, len(matches), criteria,
		)
	}
}
