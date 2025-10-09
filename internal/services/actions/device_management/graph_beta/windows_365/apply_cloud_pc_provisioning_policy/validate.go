package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

func (a *ApplyCloudPcProvisioningPolicyAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ApplyCloudPcProvisioningPolicyActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
