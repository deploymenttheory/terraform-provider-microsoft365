package graphBetaNetworkDeviceForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// validateNoUserAssignments rejects user_ids, which Graph refuses on device
// profiles ("Only device and group assignments are allowed").
func validateNoUserAssignments(ctx context.Context, rules types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	if rules.IsNull() || rules.IsUnknown() {
		return diags
	}
	var m sharedmodels.AssignmentRulesModel
	if d := rules.As(ctx, &m, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true}); d.HasError() {
		return d
	}
	if !m.UserIDs.IsNull() && !m.UserIDs.IsUnknown() {
		diags.AddAttributeError(
			path.Root("assignment_rules").AtName("user_ids"),
			"User assignments are not supported",
			"Pre-logon (device) forwarding profiles can be assigned to devices and groups only.",
		)
	}
	return diags
}
