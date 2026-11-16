package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigValidators returns the cross-attribute validators for the resource
func (r *MobilityManagementPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		scopeValidator{},
	}
}

// scopeValidator checks that applies_to, included_groups and the MDM-only setting are consistent:
//   - included_groups can only be set when applies_to is selected
//   - applies_to selected requires included_groups, because selected is applied by adding group references
//   - is_mdm_enrollment_during_registration_disabled is only valid for device policies
type scopeValidator struct{}

var _ resource.ConfigValidator = scopeValidator{}

func (v scopeValidator) Description(ctx context.Context) string {
	return "applies_to must be selected when included_groups is set, selected requires included_groups, and is_mdm_enrollment_during_registration_disabled is only valid for device policies"
}

func (v scopeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v scopeValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config MobilityManagementPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appliesToKnown := isConfigured(config.AppliesTo)
	appliesTo := config.AppliesTo.ValueString()

	if !config.IncludedGroups.IsNull() && !config.IncludedGroups.IsUnknown() && (config.AppliesTo.IsNull() || (appliesToKnown && appliesTo != appliesToSelected)) {
		resp.Diagnostics.AddAttributeError(
			path.Root("included_groups"),
			"Invalid included_groups",
			"included_groups can only be set when applies_to is \"selected\".",
		)
	}

	if appliesToKnown && appliesTo == appliesToSelected && config.IncludedGroups.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("included_groups"),
			"Missing included_groups",
			"applies_to \"selected\" requires at least one group in included_groups. Microsoft Graph applies the selected scope by adding group references.",
		)
	}

	if isConfigured(config.PolicyType) && config.PolicyType.ValueString() == policyTypeApp && !config.IsMdmEnrollmentDuringRegistrationDisabled.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("is_mdm_enrollment_during_registration_disabled"),
			"Invalid is_mdm_enrollment_during_registration_disabled",
			"is_mdm_enrollment_during_registration_disabled is only valid when policy_type is \"device\".",
		)
	}
}
