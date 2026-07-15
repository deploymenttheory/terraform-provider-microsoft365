package graphBetaNetworkContentPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkContentPolicyResourceModel represents a Global Secure Access content policy.
type NetworkContentPolicyResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Description   types.String   `tfsdk:"description"`
	DefaultAction types.String   `tfsdk:"default_action"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}
