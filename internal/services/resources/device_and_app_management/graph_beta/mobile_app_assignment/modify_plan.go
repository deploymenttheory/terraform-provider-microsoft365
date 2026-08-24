package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// intentRequired is the only install intent for which the Intune service accepts the
// isRemovable assignment setting. Any other intent is rejected with
// HTTP 400 "IsRemovable setting is only supported for Required intent."
const intentRequired = "required"

// appleIsRemovablePaths are the Apple settings blocks carrying an is_removable attribute.
var appleIsRemovablePaths = []path.Path{
	path.Root("settings").AtName("ios_lob").AtName("is_removable"),
	path.Root("settings").AtName("ios_store").AtName("is_removable"),
	path.Root("settings").AtName("ios_vpp").AtName("is_removable"),
}

// ModifyPlan reconciles the is_removable assignment setting with the install intent.
//
// The ios_lob, ios_store and ios_vpp is_removable attributes are Optional+Computed with a
// static default, so the framework materialises a value into the plan even when the
// practitioner omits the attribute entirely. That value is then sent to Graph on every
// request, and the service rejects it for any non-required intent, making it impossible to
// create an assignment for an Apple app with any other intent.
//
// Three cases are handled:
//
//   - intent is known and "required": the plan is left untouched, preserving the existing
//     default behaviour exactly.
//   - intent is known and not "required": a defaulted value is nulled, so the null-skipping
//     converters in construct.go omit the field. An explicitly configured value is an error;
//     ValidateConfig normally catches it first, but cannot when intent was unknown then.
//   - intent is unknown: a defaulted value is marked unknown. Terraform requires a value that
//     is known in the initial plan to be identical in the final plan, so leaving the default
//     known here and nulling it once intent resolves would produce an inconsistent plan.
//     Planning it as unknown lets the final plan resolve to either the default or null.
//
// Attributes are read and written by path rather than by decoding the whole plan, so that a
// configuration deriving an entire settings or target object from an unknown value is not
// rejected with a value conversion error.
func (r *MobileAppAssignmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Nothing to do when the resource is being destroyed.
	if req.Plan.Raw.IsNull() {
		return
	}

	var intent types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("intent"), &intent)...)
	if resp.Diagnostics.HasError() {
		return
	}

	intentKnown := !intent.IsNull() && !intent.IsUnknown()
	if intentKnown && intent.ValueString() == intentRequired {
		return
	}

	for _, attributePath := range appleIsRemovablePaths {
		var planned types.Bool

		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, attributePath, &planned)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// A null planned value is already omitted from the request. An unknown one is
		// resolved on a later plan, once intent is known.
		if planned.IsNull() || planned.IsUnknown() {
			continue
		}

		var configured types.Bool

		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, attributePath, &configured)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Explicitly configured by the practitioner rather than materialised from the
		// schema default, so it must not be silently discarded.
		if !configured.IsNull() {
			if intentKnown {
				addIsRemovableIntentError(&resp.Diagnostics, attributePath, intent.ValueString())
				return
			}
			continue
		}

		if !intentKnown {
			tflog.Debug(ctx, "Planning is_removable as unknown: intent is not yet known")

			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, attributePath, types.BoolUnknown())...)
			if resp.Diagnostics.HasError() {
				return
			}

			continue
		}

		tflog.Debug(ctx, "Removing is_removable from plan: only supported when intent is 'required'")

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, attributePath, types.BoolNull())...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
