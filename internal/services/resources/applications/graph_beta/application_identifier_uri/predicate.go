package graphBetaApplicationIdentifierUri

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// applicationIdentifierUriConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the identifier URI resource state has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate checks the full read state against the expected state captured at write time:
//
//   - id must be non-empty — the composite applicationId/identifierUri key; confirms Read
//     successfully located the URI in the application's identifierUris list.
//   - application_id must match expected — confirms the read is for the correct application.
//   - identifier_uri must match expected — confirms the URI propagated to the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func applicationIdentifierUriConsistencyPredicate(expected *ApplicationIdentifierUriResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ApplicationIdentifierUriResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.Id.IsNull() || actual.Id.IsUnknown() || actual.Id.ValueString() == "" {
			return false
		}

		if actual.ApplicationID.ValueString() != expected.ApplicationID.ValueString() {
			return false
		}

		if actual.IdentifierUri.ValueString() != expected.IdentifierUri.ValueString() {
			return false
		}

		return true
	}
}
