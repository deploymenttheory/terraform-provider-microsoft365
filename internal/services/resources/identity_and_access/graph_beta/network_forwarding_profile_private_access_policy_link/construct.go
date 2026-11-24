package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

const (
	forwardingPolicyLinkODataType = "#microsoft.graph.networkaccess.forwardingPolicyLink"
	forwardingPolicyODataType     = "#microsoft.graph.networkaccess.forwardingPolicy"
)

// constructResource builds the raw JSON PATCH body that adds the configured policy link to its forwarding profile
// through policies@delta.
func constructResource(plan *NetworkForwardingProfilePrivateAccessPolicyLinkResourceModel) map[string]any {
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
	return map[string]any{"policies@delta": []any{link}}
}
