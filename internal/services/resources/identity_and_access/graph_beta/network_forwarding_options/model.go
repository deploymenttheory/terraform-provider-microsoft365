package graphBetaNetworkForwardingOptions

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkForwardingOptionsResourceModel holds the tenant's managed DNS forwarding option.
type NetworkForwardingOptionsResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	SkipDNSLookupState types.String   `tfsdk:"skip_dns_lookup_state"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
