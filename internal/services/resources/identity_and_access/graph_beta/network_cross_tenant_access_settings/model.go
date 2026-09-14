package graphBetaNetworkCrossTenantAccessSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkCrossTenantAccessSettingsResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	NetworkPacketTaggingStatus types.String   `tfsdk:"network_packet_tagging_status"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}
