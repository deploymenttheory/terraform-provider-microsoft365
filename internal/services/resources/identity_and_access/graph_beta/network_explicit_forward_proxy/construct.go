package graphBetaNetworkExplicitForwardProxy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func proxyPatch(
	plan NetworkExplicitForwardProxyResourceModel,
	state *NetworkExplicitForwardProxyResourceModel,
) (map[string]any, error) {
	p := plan.InternetAccess
	if p == nil || p.SourceIPSessionAffinityOptions.IsNull() ||
		p.SourceIPSessionAffinityOptions.IsUnknown() {
		return nil, errUnknownProxySettings
	}
	for _, value := range []types.Bool{p.IsEnabled, p.IsSourceIPSessionAffinityEnabled} {
		if value.IsNull() || value.IsUnknown() {
			return nil, errUnknownProxySettings
		}
	}
	if !validAffinity(
		p.IsSourceIPSessionAffinityEnabled.ValueBool(),
		p.SourceIPSessionAffinityOptions.ValueString(),
	) {
		return nil, errUnknownProxySettings
	}
	var previous *InternetAccessResourceModel
	if state != nil {
		previous = state.InternetAccess
	}
	access := map[string]any{}
	if previous == nil || !p.IsEnabled.Equal(previous.IsEnabled) {
		access["isEnabled"] = p.IsEnabled.ValueBool()
	}
	if previous == nil ||
		!p.IsSourceIPSessionAffinityEnabled.Equal(previous.IsSourceIPSessionAffinityEnabled) ||
		!p.SourceIPSessionAffinityOptions.Equal(previous.SourceIPSessionAffinityOptions) {
		access["isSourceIpSessionAffinityEnabled"] = p.IsSourceIPSessionAffinityEnabled.ValueBool()
		if p.IsSourceIPSessionAffinityEnabled.ValueBool() {
			access["sourceIpSessionAffinityOptions"] = p.SourceIPSessionAffinityOptions.ValueString()
		}
		// Disabling affinity normalizes options to none in Graph.
	}
	if len(access) == 0 {
		return map[string]any{}, nil
	}
	return map[string]any{"internetAccess": access}, nil
}
