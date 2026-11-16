package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// groupReferenceBaseURL is the @odata.id prefix used when adding a group reference to includedGroups
const groupReferenceBaseURL = "https://graph.microsoft.com/beta/groups/"

var guidPattern = regexp.MustCompile(constants.GuidRegex)

// errInvalidResourceID is returned when an ID is not {policy_type}/{policy_id}
var errInvalidResourceID = errors.New("expected an ID in the format {policy_type}/{policy_id}, where policy_type is device or app and policy_id is a GUID")

// errReadIncludedGroups is returned when included_groups cannot be read from the plan
var errReadIncludedGroups = errors.New("failed to read included_groups from the plan")

// errPolicyBodyMismatch is returned when a request body type does not match the policy type
var errPolicyBodyMismatch = errors.New("request body type does not match the policy type")

// buildResourceID returns the composite Terraform ID {policy_type}/{policy_id}
func buildResourceID(policyType, policyID string) string {
	return policyType + "/" + policyID
}

// parseResourceID splits the composite Terraform ID {policy_type}/{policy_id}
func parseResourceID(id string) (string, string, error) {
	policyType, policyID, found := strings.Cut(id, "/")
	if !found || (policyType != policyTypeDevice && policyType != policyTypeApp) || !guidPattern.MatchString(policyID) {
		return "", "", fmt.Errorf("%w, got %q", errInvalidResourceID, id)
	}
	return policyType, policyID, nil
}

// groupReferenceURL returns the @odata.id value for a group reference
func groupReferenceURL(groupID string) string {
	return groupReferenceBaseURL + groupID
}

// diffGroupIDs returns the group IDs to add and remove so that current matches desired.
// Results are sorted so request order is deterministic.
func diffGroupIDs(current, desired []string) ([]string, []string) {
	currentSet := make(map[string]bool, len(current))
	for _, id := range current {
		currentSet[strings.ToLower(id)] = true
	}
	desiredSet := make(map[string]bool, len(desired))
	for _, id := range desired {
		desiredSet[strings.ToLower(id)] = true
	}

	toAdd := make([]string, 0)
	for _, id := range desired {
		if !currentSet[strings.ToLower(id)] {
			toAdd = append(toAdd, id)
		}
	}
	toRemove := make([]string, 0)
	for _, id := range current {
		if !desiredSet[strings.ToLower(id)] {
			toRemove = append(toRemove, id)
		}
	}

	slices.Sort(toAdd)
	slices.Sort(toRemove)

	return toAdd, toRemove
}
