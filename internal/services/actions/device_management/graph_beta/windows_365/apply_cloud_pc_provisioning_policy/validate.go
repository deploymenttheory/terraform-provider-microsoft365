package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

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

	// If reserve_percentage is set, validate that the policy is a Frontline type
	if !data.ReservePercentage.IsNull() && !data.ReservePercentage.IsUnknown() {
		if data.ProvisioningPolicyID.IsNull() || data.ProvisioningPolicyID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("reserve_percentage"),
				"Missing Required Configuration",
				"reserve_percentage can only be used when provisioning_policy_id is specified.",
			)
			return
		}

		// Fetch the provisioning policy to check its type
		policyID := data.ProvisioningPolicyID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Validating provisioning policy type for policy ID: %s", policyID))

		policy, err := a.client.
			DeviceManagement().
			VirtualEndpoint().
			ProvisioningPolicies().
			ByCloudPcProvisioningPolicyId(policyID).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("reserve_percentage"),
				"Unable to Validate Policy Type",
				fmt.Sprintf("Could not fetch provisioning policy to validate type. "+
					"Ensure the policy exists and you have permission to read it. "+
					"Error: %s", err.Error()),
			)
			return
		}

		if policy == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("provisioning_policy_id"),
				"Policy Not Found",
				fmt.Sprintf("Provisioning policy with ID %s was not found.", policyID),
			)
			return
		}

		// Check if the policy is a Frontline type (shared, sharedByUser, or sharedByEntraGroup)
		provisioningType := policy.GetProvisioningType()
		if provisioningType == nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("reserve_percentage"),
				"Unable to Determine Policy Type",
				"Could not determine the provisioning type of the policy. "+
					"reserve_percentage is only applicable for Frontline shared provisioning policies.",
			)
			return
		}

		provisioningTypeStr := provisioningType.String()
		isFrontline := provisioningTypeStr == "shared" ||
			provisioningTypeStr == "sharedByUser" ||
			provisioningTypeStr == "sharedByEntraGroup"

		if !isFrontline {
			resp.Diagnostics.AddAttributeError(
				path.Root("reserve_percentage"),
				"Invalid Configuration",
				fmt.Sprintf("reserve_percentage can only be used with Frontline shared provisioning policies. "+
					"The policy %s has provisioning type '%s', which is not a Frontline type. "+
					"Valid Frontline types are: shared, sharedByUser, sharedByEntraGroup.",
					policyID, provisioningTypeStr),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Validation successful: Policy %s is Frontline type '%s'", policyID, provisioningTypeStr))
	}
}
