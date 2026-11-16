package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// lifecycleWorkflowConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// accepts a read only once it reflects the preceding write (ADR-002).
//
// Probing the tenant showed create, PATCH and createNewVersion are readable immediately, so the
// predicate normally passes on the first read; it guards against a replica serving the pre-update
// workflow.
//
// prevModified is last_modified_date_time captured from state before the write:
//   - For Create: pass types.StringNull(); the predicate accepts once the workflow has an ID and a
//     last_modified_date_time.
//   - For Update: pass state.LastModifiedDateTime; the predicate accepts once the timestamp differs,
//     confirming the read returned the PATCHed or newly versioned workflow.
//
// prevVersion is the version captured before a createNewVersion call (types.Int64Null() otherwise);
// when set, the read must report a higher version.
func lifecycleWorkflowConsistencyPredicate(prevModified types.String, prevVersion types.Int64) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual LifecycleWorkflowResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.LastModifiedDateTime.IsNull() || actual.LastModifiedDateTime.IsUnknown() || actual.LastModifiedDateTime.ValueString() == "" {
			return false
		}

		if !prevModified.IsNull() && !prevModified.IsUnknown() && prevModified.ValueString() != "" {
			if actual.LastModifiedDateTime.ValueString() == prevModified.ValueString() {
				return false
			}
		}

		if !prevVersion.IsNull() && !prevVersion.IsUnknown() {
			if actual.Version.IsNull() || actual.Version.IsUnknown() || actual.Version.ValueInt64() <= prevVersion.ValueInt64() {
				return false
			}
		}

		return true
	}
}
