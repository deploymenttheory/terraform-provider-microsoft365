package graphBetaNetworkExplicitForwardProxy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValidateConfig validates the relationship between the affinity flag and options.
func (r *NetworkExplicitForwardProxyResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var enabled types.Bool
	var options types.String
	resp.Diagnostics.Append(
		req.Config.GetAttribute(
			ctx,
			path.Root("internet_access").AtName("is_source_ip_session_affinity_enabled"),
			&enabled,
		)...)
	resp.Diagnostics.Append(
		req.Config.GetAttribute(
			ctx,
			path.Root("internet_access").AtName("source_ip_session_affinity_options"),
			&options,
		)...)
	if resp.Diagnostics.HasError() || enabled.IsNull() || enabled.IsUnknown() || options.IsNull() ||
		options.IsUnknown() {
		return
	}
	if !validAffinity(enabled.ValueBool(), options.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("internet_access").AtName("source_ip_session_affinity_options"),
			"Inconsistent session affinity settings",
			"Use none when is_source_ip_session_affinity_enabled is false. When true, select useSessionId, useHttpHeader, or useSessionId,useHttpHeader.",
		)
	}
}

func validAffinity(enabled bool, options string) bool {
	switch options {
	case "none":
		return !enabled
	case "useSessionId", "useHttpHeader", "useSessionId,useHttpHeader":
		return enabled
	default:
		return false
	}
}
