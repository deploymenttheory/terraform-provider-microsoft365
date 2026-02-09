package graphBetaConditionalAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConditionalAccessPolicyListConfigModel represents the configuration for listing Conditional Access policies
type ConditionalAccessPolicyListConfigModel struct {
	DisplayNameFilter types.String `tfsdk:"display_name_filter"`
	StateFilter       types.String `tfsdk:"state_filter"`
	ODataFilter       types.String `tfsdk:"odata_filter"`
}
