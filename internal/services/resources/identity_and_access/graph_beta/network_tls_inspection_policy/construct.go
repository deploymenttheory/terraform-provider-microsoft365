package graphBetaNetworkTLSInspectionPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

// constructResource uses the generated TLS inspection policy model. defaultAction
// remains in additional data because it is currently missing from Graph metadata.
func constructResource(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyResourceModel,
) (models.TlsInspectionPolicyable, error) {
	body := models.NewTlsInspectionPolicy()
	body.SetName(data.Name.ValueStringPointer())
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		body.SetDescription(data.Description.ValueStringPointer())
	}
	settings := models.NewTlsInspectionPolicySettings()
	settings.GetAdditionalData()["defaultAction"] = data.DefaultAction.ValueString()
	body.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		body,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return body, nil
}

// constructUpdateResource sends only changed fields. A nil description is placed
// in additional data because generated string setters omit nil values.
func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkTLSInspectionPolicyResourceModel,
) (models.TlsInspectionPolicyable, error) {
	body := models.NewTlsInspectionPolicy()
	if !plan.Name.Equal(state.Name) {
		body.SetName(plan.Name.ValueStringPointer())
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			body.GetAdditionalData()["description"] = nil
		} else {
			body.SetDescription(plan.Description.ValueStringPointer())
		}
	}
	if !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := models.NewTlsInspectionPolicySettings()
		settings.GetAdditionalData()["defaultAction"] = plan.DefaultAction.ValueString()
		body.SetSettings(settings)
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		body,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return body, nil
}

func hasUpdateChanges(plan, state *NetworkTLSInspectionPolicyResourceModel) bool {
	return !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) ||
		!plan.DefaultAction.Equal(state.DefaultAction)
}
