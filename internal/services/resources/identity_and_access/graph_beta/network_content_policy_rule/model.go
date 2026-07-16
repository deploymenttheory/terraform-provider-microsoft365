package graphBetaNetworkContentPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkContentPolicyRuleResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	ContentPolicyID  types.String   `tfsdk:"content_policy_id"`
	Name             types.String   `tfsdk:"name"`
	Description      types.String   `tfsdk:"description"`
	Action           types.String   `tfsdk:"action"`
	Priority         types.Int64    `tfsdk:"priority"`
	Status           types.String   `tfsdk:"status"`
	Activities       types.Set      `tfsdk:"activities"`
	ContentTypes     types.Set      `tfsdk:"content_types"`
	TextContentTypes types.Set      `tfsdk:"text_content_types"`
	Destinations     types.List     `tfsdk:"destinations"`
	SessionTypes     types.Set      `tfsdk:"session_types"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type ContentPolicyRuleDestinationModel struct {
	Type   types.String `tfsdk:"type"`
	Values types.Set    `tfsdk:"values"`
}
