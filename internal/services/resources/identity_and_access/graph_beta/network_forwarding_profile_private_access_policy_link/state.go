package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteStateToTerraform maps the policy link onto state. It keeps the configured casing for IDs that Graph echoes
// back case-insensitively.
func MapRemoteStateToTerraform(state *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel, link *policyLink) {
	if !strings.EqualFold(state.PolicyLinkID.ValueString(), link.ID) {
		state.PolicyLinkID = types.StringValue(link.ID)
	}
	state.ID = types.StringNull()
	if profileID, linkID := state.ForwardingProfileID.ValueString(), state.PolicyLinkID.ValueString(); profileID != "" && linkID != "" {
		state.ID = types.StringValue(profileID + "/" + linkID)
	}
	if link.Policy != nil && !strings.EqualFold(state.ForwardingPolicyID.ValueString(), link.Policy.ID) {
		state.ForwardingPolicyID = types.StringValue(link.Policy.ID)
	}
	if link.ConfigType != "" {
		state.ConfigType = types.StringValue(link.ConfigType)
	}
	if link.State != "" {
		state.State = types.StringValue(strings.ToLower(link.State))
	}
	if link.Priority != nil {
		state.Priority = types.Int64Value(*link.Priority)
	} else {
		state.Priority = types.Int64Null()
	}
}

// findPolicyLinkByIDOrPolicyID returns the link with linkID, or else the link to policyID, because a re-added link gets
// a new ID. An empty linkID or policyID skips that match.
func findPolicyLinkByIDOrPolicyID(links []policyLink, linkID, policyID string) *policyLink {
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
