package graphBetaNetworkMCPPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkMCPPolicyResourceModel represents a Global Secure Access MCP policy.
type NetworkMCPPolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	DefaultAction        types.String   `tfsdk:"default_action"`
	Version              types.String   `tfsdk:"version"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
