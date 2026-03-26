package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// agentIdentityBlueprintServicePrincipalConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the agent identity blueprint service principal write has propagated
// across Microsoft Entra replicas before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// app_id is an immutable key field — it identifies the application registration that backs this
// service principal. Tags is the sole mutable user-specified field.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func agentIdentityBlueprintServicePrincipalConsistencyPredicate(expected *AgentIdentityBlueprintServicePrincipalResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AgentIdentityBlueprintServicePrincipalResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.AppId.ValueString() != expected.AppId.ValueString() {
			return false
		}

		// Tags is the sole mutable user-specified field on this resource.
		if !expected.Tags.IsNull() && !expected.Tags.IsUnknown() {
			if !actual.Tags.Equal(expected.Tags) {
				return false
			}
		}

		return true
	}
}
