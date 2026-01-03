// REF: https://learn.microsoft.com/en-us/graph/api/cloudpcprovisioningpolicy-apply?view=graph-rest-beta
package graphBetaApplyCloudPcProvisioningPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplyCloudPcProvisioningPolicyActionModel struct {
	ProvisioningPolicyID types.String   `tfsdk:"provisioning_policy_id"`
	PolicySettings       types.String   `tfsdk:"policy_settings"`
	ReservePercentage    types.Int32    `tfsdk:"reserve_percentage"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
