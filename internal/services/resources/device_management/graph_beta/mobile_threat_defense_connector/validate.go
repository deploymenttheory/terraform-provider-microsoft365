package graphBetaMobileThreatDefenseConnector

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const defenderAttachOnlyMessage = "Only the Microsoft Defender for Endpoint connector (%s) supports this setting. Graph ignores it for partner %s."

// errConnectorAlreadyExists is returned when create targets a partner that already has a connector.
var errConnectorAlreadyExists = errors.New("mobile threat defense connector already exists")

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

// validateRequest refuses to create a connector for a partner that already has one.
// Graph's POST is an upsert keyed on the partner ID, so posting to an existing connector would silently replace its settings.
// Only runs during create operations when isCreate is true.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *MobileThreatDefenseConnectorResourceModel, isCreate bool) error {
	if !isCreate {
		tflog.Debug(ctx, "Skipping existing connector check for update operation", map[string]any{"partnerId": data.ID.ValueString()})
		return nil
	}

	partnerId := data.ID.ValueString()

	tflog.Debug(ctx, "Checking whether the mobile threat defense connector already exists", map[string]any{"partnerId": partnerId})

	_, err := client.DeviceManagement().MobileThreatDefenseConnectors().ByMobileThreatDefenseConnectorId(partnerId).Get(ctx, nil)
	if err == nil {
		return fmt.Errorf("%w for partner %s: import it instead with 'terraform import %s.<name> %s'", errConnectorAlreadyExists, partnerId, ResourceName, partnerId)
	}

	if kiotaerrors.GraphError(ctx, err).StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, "No existing mobile threat defense connector found", map[string]any{"partnerId": partnerId})
		return nil
	}

	return fmt.Errorf("could not check whether mobile threat defense connector %s exists: %w", partnerId, err)
}
