package graphBetaMacOSDepEnrollmentProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigValidators implements resource.ResourceWithConfigValidators. It encodes the
// Intune cross-field business rule that a userless or auto-advancing macOS DEP profile
// must be mandatory, so users get a clear plan-time error instead of a Graph 400 at apply.
func (r *MacOSDepEnrollmentProfileResource) ConfigValidators(
	ctx context.Context,
) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		requireMandatoryWhenUserlessOrAutoAdvance{},
	}
}

// requireMandatoryWhenUserlessOrAutoAdvance enforces:
// if requires_user_authentication = false OR auto_advance_setup_enabled = true,
// then is_mandatory must be set to true.
type requireMandatoryWhenUserlessOrAutoAdvance struct{}

func (v requireMandatoryWhenUserlessOrAutoAdvance) Description(_ context.Context) string {
	return "is_mandatory must be true when requires_user_authentication is false or auto_advance_setup_enabled is true"
}

func (v requireMandatoryWhenUserlessOrAutoAdvance) MarkdownDescription(ctx context.Context) string {
	return "`is_mandatory` must be `true` when `requires_user_authentication` is `false` or `auto_advance_setup_enabled` is `true`"
}

func (v requireMandatoryWhenUserlessOrAutoAdvance) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if req.Config.Raw.IsNull() {
		return
	}

	var requiresUserAuth, autoAdvance, isMandatory types.Bool
	resp.Diagnostics.Append(
		req.Config.GetAttribute(
			ctx,
			path.Root("requires_user_authentication"),
			&requiresUserAuth,
		)...)
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("auto_advance_setup_enabled"), &autoAdvance)...)
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("is_mandatory"), &isMandatory)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip while values are still unknown (e.g. interpolated at plan time).
	if requiresUserAuth.IsUnknown() || autoAdvance.IsUnknown() || isMandatory.IsUnknown() {
		return
	}

	userless := !requiresUserAuth.IsNull() && !requiresUserAuth.ValueBool()
	autoAdvancing := !autoAdvance.IsNull() && autoAdvance.ValueBool()
	if !userless && !autoAdvancing {
		return
	}

	// is_mandatory must be explicitly true. Unset (null) or false both violate the rule
	// because Graph rejects mandatory=false for these profiles.
	if isMandatory.IsNull() || !isMandatory.ValueBool() {
		var trigger string
		switch {
		case userless && autoAdvancing:
			trigger = "requires_user_authentication is false and auto_advance_setup_enabled is true"
		case userless:
			trigger = "requires_user_authentication is false (userless enrollment)"
		default:
			trigger = "auto_advance_setup_enabled is true"
		}
		resp.Diagnostics.AddAttributeError(
			path.Root("is_mandatory"),
			"is_mandatory must be true",
			"This macOS DEP enrollment profile requires is_mandatory = true because "+trigger+". "+
				"Microsoft Graph rejects a non-mandatory profile in this configuration. Set is_mandatory = true.",
		)
	}
}
