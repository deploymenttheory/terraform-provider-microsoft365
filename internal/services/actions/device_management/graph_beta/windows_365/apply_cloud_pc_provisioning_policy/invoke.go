package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ApplyCloudPcProvisioningPolicyAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ApplyCloudPcProvisioningPolicyActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.ProvisioningPolicyID.ValueString()
	policySettings := data.PolicySettings.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing action %s, applying policy settings '%s' for policy ID: %s", ActionName, policySettings, policyID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Applying %s settings to provisioning policy %s...", policySettings, policyID),
	})

	requestBody, err := constructRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing request",
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Action", a.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully applied %s settings to provisioning policy %s", policySettings, policyID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Policy settings applied successfully to %s", policyID),
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}
