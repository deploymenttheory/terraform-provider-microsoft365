package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

const (
	forwardingPolicyLinkODataType = "#microsoft.graph.networkaccess.forwardingPolicyLink"
	forwardingPolicyODataType     = "#microsoft.graph.networkaccess.forwardingPolicy"
	// A matching read can be followed by a stale one.
	stableReads = 3
)

// profileTemplate returns the API path for a user or device forwarding profile.
func profileTemplate(profileType string) string {
	if profileType == profileTypeDevice {
		return "/networkAccess/deviceForwardingProfiles/{forwardingProfileId}"
	}
	return "/networkAccess/forwardingProfiles/{forwardingProfileId}"
}

// constructLink builds the forwarding policy link payload.
func constructLink(plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) map[string]any {
	link := map[string]any{
		"@odata.type": forwardingPolicyLinkODataType,
		"state":       plan.State.ValueString(),
		"configType":  plan.ConfigType.ValueString(),
		"policy": map[string]any{
			"@odata.type": forwardingPolicyODataType,
			"id":          plan.ForwardingPolicyID.ValueString(),
		},
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		link["priority"] = plan.Priority.ValueInt64()
	}
	return link
}

// constructPolicyDelta builds a policy-link delta without changing other links.
func constructPolicyDelta(delta any) map[string]any {
	return map[string]any{"policies@delta": []any{delta}}
}

// constructRemovedLink preserves the Graph-required ordering of removal metadata.
func constructRemovedLink(linkID string) removedLink {
	return removedLink{Removed: removedReason{Reason: "deleted"}, ODataType: forwardingPolicyLinkODataType, ID: linkID}
}
