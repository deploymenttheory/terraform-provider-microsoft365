package graphBetaMobileThreatDefenseConnector

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const defenderAttachOnlyMessage = "Only the Microsoft Defender for Endpoint connector (%s) supports this setting. Graph ignores it for partner %s."

// ValidateConfig rejects settings that Graph silently ignores for the configured partner,
// which would otherwise surface as an inconsistent result after apply.
func (r *MobileThreatDefenseConnectorResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data MobileThreatDefenseConnectorResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() || data.ID.IsUnknown() || data.MicrosoftDefenderForEndpointAttachEnabled.IsNull() || data.MicrosoftDefenderForEndpointAttachEnabled.IsUnknown() {
		return
	}

	if data.MicrosoftDefenderForEndpointAttachEnabled.ValueBool() && data.ID.ValueString() != MicrosoftDefenderForEndpointPartnerId {
		detail := fmt.Sprintf(defenderAttachOnlyMessage, MicrosoftDefenderForEndpointPartnerId, data.ID.ValueString())
		resp.Diagnostics.AddAttributeError(path.Root("microsoft_defender_for_endpoint_attach_enabled"), "Unsupported setting for partner", detail)
	}
}
