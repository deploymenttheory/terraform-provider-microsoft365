package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// agentIdentityBlueprintIdentifierUriConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the identifier URI resource state has propagated across Microsoft
// Entra replicas before accepting the read as authoritative.
//
// The predicate checks the full read state against the expected state captured at write time:
//
//   - id must be non-empty — confirms the resource was found by Read.
//   - identifier_uri must match expected — confirms the URI propagated to the responding replica.
//   - blueprint_id must match expected — confirms the read is for the correct application.
//   - if a scope was configured, scope.id must be non-empty — a Computed field assigned by Entra
//     when the OAuth2 permission scope is created.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func agentIdentityBlueprintIdentifierUriConsistencyPredicate(expected *AgentIdentityBlueprintIdentifierUriResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AgentIdentityBlueprintIdentifierUriResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.BlueprintID.ValueString() != expected.BlueprintID.ValueString() {
			return false
		}

		if actual.IdentifierUri.ValueString() != expected.IdentifierUri.ValueString() {
			return false
		}

		// If a scope was configured, scope.id is Computed (assigned by Entra) and must be present.
		if expected.Scope != nil {
			if actual.Scope == nil {
				return false
			}
			if actual.Scope.ID.IsNull() || actual.Scope.ID.IsUnknown() || actual.Scope.ID.ValueString() == "" {
				return false
			}
		}

		return true
	}
}
