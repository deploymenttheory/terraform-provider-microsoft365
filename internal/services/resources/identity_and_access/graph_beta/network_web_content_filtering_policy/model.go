package graphBetaNetworkWebContentFilteringPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkWebContentFilteringPolicyResourceModel represents the Terraform schema
// for Global Secure Access web content filtering policies.
type NetworkWebContentFilteringPolicyResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Description   types.String   `tfsdk:"description"`
	DefaultAction types.String   `tfsdk:"default_action"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}
