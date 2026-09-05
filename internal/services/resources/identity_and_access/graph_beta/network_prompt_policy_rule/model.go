package graphBetaNetworkPromptPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkPromptPolicyRuleResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	PromptPolicyID      types.String   `tfsdk:"prompt_policy_id"`
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	Action              types.String   `tfsdk:"action"`
	Priority            types.Int32    `tfsdk:"priority"`
	Enabled             types.Bool     `tfsdk:"enabled"`
	Status              types.String   `tfsdk:"status"`
	ConversationSchemes types.List     `tfsdk:"conversation_schemes"`
	PromptLogging       types.String   `tfsdk:"prompt_logging"`
	ScanResult          types.String   `tfsdk:"scan_result"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

type ConversationSchemeModel struct {
	Type       types.String `tfsdk:"type"`
	URL        types.String `tfsdk:"url"`
	JSONPath   types.String `tfsdk:"json_path"`
	SchemeName types.String `tfsdk:"scheme_name"`
}

type PromptPolicyRuleIdentity struct {
	ID             string `tfsdk:"id"`
	PromptPolicyID string `tfsdk:"prompt_policy_id"`
}
