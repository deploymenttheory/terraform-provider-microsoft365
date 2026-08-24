package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValidateConfig rejects configurations that explicitly set is_removable on an Apple settings
// block alongside an install intent other than "required".
//
// The Intune service only supports the isRemovable setting for required intent and returns
// HTTP 400 "IsRemovable setting is only supported for Required intent." otherwise. Reporting
// it at validation time gives the practitioner an actionable error rather than a failure
// part-way through an apply.
//
// When intent is not yet known this cannot be evaluated; ModifyPlan repeats the check once
// intent resolves.
//
// Attributes are read by path rather than by decoding the whole configuration, so that a
// configuration deriving an entire settings or target object from an unknown value is not
// rejected with a value conversion error.
func (r *MobileAppAssignmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var intent types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("intent"), &intent)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if intent.IsNull() || intent.IsUnknown() || intent.ValueString() == intentRequired {
		return
	}

	for _, attributePath := range appleIsRemovablePaths {
		var isRemovable types.Bool

		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, attributePath, &isRemovable)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if isRemovable.IsNull() || isRemovable.IsUnknown() {
			continue
		}

		addIsRemovableIntentError(&resp.Diagnostics, attributePath, intent.ValueString())
	}
}

// addIsRemovableIntentError appends the diagnostic raised when is_removable is configured
// alongside an install intent that the service does not support it for.
func addIsRemovableIntentError(diagnostics *diag.Diagnostics, attributePath path.Path, intent string) {
	diagnostics.AddAttributeError(
		attributePath,
		"Invalid Attribute Combination",
		fmt.Sprintf(
			"is_removable cannot be set when intent is %q. The Intune service only supports "+
				"this setting when intent is %q. Remove the attribute, or change the intent.",
			intent, intentRequired,
		),
	)
}
