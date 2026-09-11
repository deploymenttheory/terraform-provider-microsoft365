package graphBetaNetworkTLSInspectionPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkTLSInspectionPolicyRuleResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	TLSInspectionPolicyID types.String   `tfsdk:"tls_inspection_policy_id"`
	Name                  types.String   `tfsdk:"name"`
	Description           types.String   `tfsdk:"description"`
	Action                types.String   `tfsdk:"action"`
	Priority              types.Int32    `tfsdk:"priority"`
	Enabled               types.Bool     `tfsdk:"enabled"`
	Status                types.String   `tfsdk:"status"`
	Destinations          types.List     `tfsdk:"destinations"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

type TLSInspectionPolicyRuleDestinationModel struct {
	Type   types.String `tfsdk:"type"`
	Values types.Set    `tfsdk:"values"`
}

type TLSInspectionPolicyRuleIdentity struct {
	ID                    string `tfsdk:"id"`
	TLSInspectionPolicyID string `tfsdk:"tls_inspection_policy_id"`
}
