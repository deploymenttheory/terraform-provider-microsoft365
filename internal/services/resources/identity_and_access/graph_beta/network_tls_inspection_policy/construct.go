package graphBetaNetworkTLSInspectionPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

// constructResource builds the TLS inspection policy body. Create passes a nil state and sends every field; Update passes
// the prior state and sends only the fields that changed. settings.defaultAction goes through additional data because the
// SDK model lacks it (see the SDK deviation in crud.go).
func constructResource(ctx context.Context, plan, state *NetworkTLSInspectionPolicyResourceModel) (models.TlsInspectionPolicyable, error) {
	body := models.NewTlsInspectionPolicy()
	if state == nil || !plan.Name.Equal(state.Name) {
		body.SetName(plan.Name.ValueStringPointer())
	}
	if state == nil {
		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			body.SetDescription(plan.Description.ValueStringPointer())
		}
	} else if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			// The generated setter omits nil strings, so an explicit JSON null clears the description through additional data.
			body.GetAdditionalData()["description"] = nil
		} else {
			body.SetDescription(plan.Description.ValueStringPointer())
		}
	}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := models.NewTlsInspectionPolicySettings()
		settings.GetAdditionalData()["defaultAction"] = plan.DefaultAction.ValueString()
		body.SetSettings(settings)
	}
	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), body); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return body, nil
}
