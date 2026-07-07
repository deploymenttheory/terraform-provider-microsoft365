package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest performs pre-request validation for Update. It rejects an attempt to flip
// is_default_policy_assignment from true to false while this policy is still the DEP token's
// current default iOS/iPadOS enrollment profile on Graph.
//
// Microsoft Graph exposes no unset/clear action for the default profile - the only way a policy
// stops being the default is another iOS/iPadOS profile being promoted via setDefaultProfile. So a
// true -> false change on the current default can never converge on its own: the PUT would
// succeed, Graph would keep the policy as the default, and the post-apply read would re-derive
// true, failing the apply with "Provider produced inconsistent result". This is checked against
// the token's live default (GET /deviceManagement/depOnboardingSettings/{id}?$expand=defaultIosEnrollmentProfile)
// rather than state, so a demote that is ordered after another policy's promotion in the same
// apply (via depends_on) validates cleanly. A live lookup is required, which is why this cannot
// be expressed as a schema or config validator.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) validateRequest(ctx context.Context, plan, state *IOSiPadOSDeviceEnrollmentPolicyResourceModel) error {
	stateIsDefault := !state.IsDefaultPolicyAssignment.IsNull() && state.IsDefaultPolicyAssignment.ValueBool()
	planUnsetsDefault := !plan.IsDefaultPolicyAssignment.IsNull() && !plan.IsDefaultPolicyAssignment.IsUnknown() &&
		!plan.IsDefaultPolicyAssignment.ValueBool()

	if !stateIsDefault || !planUnsetsDefault {
		return nil
	}

	depOnboardingSettingsId := state.DepOnboardingSettingsId.ValueString()
	policyId := state.ID.ValueString()
	if depOnboardingSettingsId == "" || policyId == "" {
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("is_default_policy_assignment true->false requested for policy %s; checking whether it is still the current default iOS/iPadOS enrollment profile for dep_onboarding_settings_id %s",
		policyId, depOnboardingSettingsId))

	isCurrentDefault, err := r.resolveIsDefaultPolicyAssignment(ctx, depOnboardingSettingsId, policyId)
	if err != nil {
		return fmt.Errorf("could not resolve the current default iOS/iPadOS enrollment profile for dep_onboarding_settings_id %s to verify that is_default_policy_assignment can be unset: %w",
			depOnboardingSettingsId, err)
	}

	if isCurrentDefault {
		return fmt.Errorf("cannot unset the default policy assignment: this policy (%s) is currently the default iOS/iPadOS enrollment profile for its DEP token "+
			"(dep_onboarding_settings_id: %s). Microsoft Graph provides no action to unset a default enrollment profile; a policy only stops being the "+
			"default when a different iOS/iPadOS enrollment policy on the same token is promoted via its own is_default_policy_assignment = true. Promote the "+
			"replacement policy first (in the same apply, add a depends_on for the replacement to this policy, or apply the promotion separately) - once "+
			"another policy is the default, this attribute refreshes to false on its own",
			policyId, depOnboardingSettingsId)
	}

	return nil
}
