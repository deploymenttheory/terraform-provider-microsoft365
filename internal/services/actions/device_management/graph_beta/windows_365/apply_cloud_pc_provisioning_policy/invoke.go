package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ApplyCloudPcProvisioningPolicyAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ApplyCloudPcProvisioningPolicyActionModel

	tflog.Debug(ctx, "Starting Cloud PC provisioning policy application", map[string]any{"action": ActionName})

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Invoke, InvokeTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	policyID := data.ProvisioningPolicyID.ValueString()
	policySettings := "region"
	if !data.PolicySettings.IsNull() && !data.PolicySettings.IsUnknown() {
		policySettings = data.PolicySettings.ValueString()
	}

	tflog.Debug(ctx, "Processing Cloud PC provisioning policy", map[string]any{
		"policy_id":       policyID,
		"policy_settings": policySettings,
	})

	validatePolicyExists := true
	if !data.ValidatePolicyExists.IsNull() && !data.ValidatePolicyExists.IsUnknown() {
		validatePolicyExists = data.ValidatePolicyExists.ValueBool()
	}

	reservePercentageSet := !data.ReservePercentage.IsNull() && !data.ReservePercentage.IsUnknown()

	if validatePolicyExists {
		tflog.Debug(ctx, "Performing policy validation via API")

		validationResult, err := validateRequest(ctx, a.client, policyID, reservePercentageSet)
		if err != nil {
			tflog.Error(ctx, "Failed to validate policy via API", map[string]any{"error": err.Error()})
			resp.Diagnostics.AddError("Policy Validation Failed", fmt.Sprintf("Failed to validate provisioning policy: %s", err.Error()))
			return
		}

		if validationResult.PolicyNotFound {
			resp.Diagnostics.AddError(
				"Provisioning Policy Not Found",
				fmt.Sprintf("Provisioning policy with ID %s does not exist or is not accessible. "+
					"Ensure the policy ID is correct and you have permission to access it.", policyID),
			)
			return
		}

		if validationResult.InvalidProvisionType {
			resp.Diagnostics.AddError(
				"Invalid Provisioning Type",
				fmt.Sprintf("reserve_percentage can only be used with Frontline shared provisioning policies. "+
					"The policy %s has provisioning type '%s', which is not a Frontline type. "+
					"Valid Frontline types are: shared, sharedByUser, sharedByEntraGroup.",
					policyID, validationResult.ProvisioningType),
			)
			return
		}

		tflog.Debug(ctx, "Policy validation completed successfully")
	} else {
		tflog.Debug(ctx, "Policy validation disabled, skipping API checks")
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Applying %s settings to provisioning policy %s...", policySettings, policyID),
	})

	requestBody, err := constructRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Constructing Request",
			fmt.Sprintf("Could not construct request for apply provisioning policy: %s", err.Error()),
		)
		return
	}

	err = a.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(policyID).
		Apply().
		Post(ctx, requestBody, nil)

	if err != nil {
		tflog.Error(ctx, "Failed to apply policy settings", map[string]any{"policy_id": policyID, "error": err.Error()})
		errors.HandleKiotaGraphError(ctx, err, resp, "Action", a.WritePermissions)
		return
	}

	tflog.Info(ctx, "Successfully applied policy settings", map[string]any{
		"policy_id":       policyID,
		"policy_settings": policySettings,
	})

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully applied %s settings to provisioning policy %s. "+
			"Changes will apply to existing Cloud PCs when they are not actively connected. "+
			"Cloud PCs are reprovisioned only when there are no active users.",
			policySettings, policyID),
	})

	tflog.Info(ctx, "Cloud PC provisioning policy application completed", map[string]any{
		"policy_id":       policyID,
		"policy_settings": policySettings,
	})
}
