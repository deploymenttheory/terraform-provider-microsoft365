package graphBetaNetworkExplicitForwardProxy

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var errUnknownProxySettings = errors.New("internet_access attributes must be known and non-null, with affinity enabled exactly when options is not none")

// affinityOptionsEnabled maps each supported sourceIpSessionAffinityOptions value to the isSourceIpSessionAffinityEnabled
// value it requires: none only when affinity is disabled, every other option only when it is enabled.
var affinityOptionsEnabled = map[string]bool{
	"none":                       false,
	"useSessionId":               true,
	"useHttpHeader":              true,
	"useSessionId,useHttpHeader": true,
}

// constructResource builds the raw JSON PATCH body for the explicit forward proxy singleton. With a nil state it sends
// every internetAccess setting; otherwise it sends only the changed ones, and an empty map when nothing changed.
func constructResource(plan NetworkExplicitForwardProxyResourceModel, state *NetworkExplicitForwardProxyResourceModel) (map[string]any, error) {
	p := plan.InternetAccess
	if p == nil || p.SourceIPSessionAffinityOptions.IsNull() || p.SourceIPSessionAffinityOptions.IsUnknown() {
		return nil, errUnknownProxySettings
	}
	for _, value := range []types.Bool{p.IsEnabled, p.IsSourceIPSessionAffinityEnabled} {
		if value.IsNull() || value.IsUnknown() {
			return nil, errUnknownProxySettings
		}
	}
	if wantEnabled, ok := affinityOptionsEnabled[p.SourceIPSessionAffinityOptions.ValueString()]; !ok || wantEnabled != p.IsSourceIPSessionAffinityEnabled.ValueBool() {
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
	affinityChanged := previous == nil || !p.IsSourceIPSessionAffinityEnabled.Equal(previous.IsSourceIPSessionAffinityEnabled)
	if affinityChanged || !p.SourceIPSessionAffinityOptions.Equal(previous.SourceIPSessionAffinityOptions) {
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
