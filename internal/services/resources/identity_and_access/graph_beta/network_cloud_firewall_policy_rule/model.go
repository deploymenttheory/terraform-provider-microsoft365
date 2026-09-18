package graphBetaNetworkCloudFirewallPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkCloudFirewallPolicyRuleResourceModel struct {
	ID           types.String   `tfsdk:"id"`
	PolicyID     types.String   `tfsdk:"policy_id"`
	Name         types.String   `tfsdk:"name"`
	Description  types.String   `tfsdk:"description"`
	Priority     types.Int32    `tfsdk:"priority"`
	Action       types.String   `tfsdk:"action"`
	Enabled      types.Bool     `tfsdk:"enabled"`
	Status       types.String   `tfsdk:"status"`
	Sources      types.Object   `tfsdk:"sources"`
	Destinations types.Object   `tfsdk:"destinations"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

type ruleIdentity struct {
	PolicyID string `tfsdk:"policy_id"`
	ID       string `tfsdk:"id"`
}

type ruleAddressModel struct {
	Type   types.String `tfsdk:"type"`
	Values types.Set    `tfsdk:"values"`
}

func addressObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "values": types.SetType{ElemType: types.StringType}}}
}

func matchingObjectType(destination bool) types.ObjectType {
	attrs := map[string]attr.Type{"addresses": types.SetType{ElemType: addressObjectType()}, "ports": types.SetType{ElemType: types.StringType}}
	if destination {
		attrs["protocols"] = types.SetType{ElemType: types.StringType}
	}
	return types.ObjectType{AttrTypes: attrs}
}
