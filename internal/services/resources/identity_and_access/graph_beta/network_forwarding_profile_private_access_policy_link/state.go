package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"strings"
)

// findLink falls back to the policy ID because a re-added link gets a new ID.
func findLink(links []policyLink, linkID, policyID string) *policyLink {
	if linkID != "" {
		for i := range links {
			if strings.EqualFold(links[i].ID, linkID) {
				return &links[i]
			}
		}
	}
	if policyID != "" {
		for i := range links {
			if links[i].Policy != nil && strings.EqualFold(links[i].Policy.ID, policyID) {
				return &links[i]
			}
		}
	}
	return nil
}
