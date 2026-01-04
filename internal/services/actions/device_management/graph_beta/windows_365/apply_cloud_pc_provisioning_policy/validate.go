package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ApplyCloudPcProvisioningPolicyAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ApplyCloudPcProvisioningPolicyActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that reserve_percentage requires provisioning_policy_id
	if !data.ReservePercentage.IsNull() && !data.ReservePercentage.IsUnknown() {
		if data.ProvisioningPolicyID.IsNull() || data.ProvisioningPolicyID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("reserve_percentage"),
				"Missing Required Configuration",
				"reserve_percentage can only be used when provisioning_policy_id is specified.",
			)
			return
		}

		// Warning about Frontline requirement
		resp.Diagnostics.AddAttributeWarning(
			path.Root("reserve_percentage"),
			"Frontline Policy Required",
			"reserve_percentage is only applicable for Frontline shared provisioning policies (shared, sharedByUser, sharedByEntraGroup). "+
				"Ensure the specified provisioning policy is a Frontline type, otherwise this parameter will be ignored.",
		)
	}

	// Informational message about policy application
	policySettings := "region"
	if !data.PolicySettings.IsNull() && !data.PolicySettings.IsUnknown() {
		policySettings = data.PolicySettings.ValueString()
	}

	resp.Diagnostics.AddAttributeWarning(
		path.Root("policy_settings"),
		"Cloud PC Provisioning Policy Application",
		"This action will apply "+policySettings+" settings to all existing Cloud PCs that were provisioned with this policy. "+
			"Cloud PCs are reprovisioned only when there are no active and connected users. "+
			"Note: Network and image changes cannot be applied retrospectively and require reprovisioning.",
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"policy_id":       data.ProvisioningPolicyID.ValueString(),
		"policy_settings": policySettings,
	})
}
