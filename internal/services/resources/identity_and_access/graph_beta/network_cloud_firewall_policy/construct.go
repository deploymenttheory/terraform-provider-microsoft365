package graphBetaNetworkCloudFirewallPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// policyRequest tracks changes to API-managed fields.
type policyRequest struct {
	*graph.CloudFirewallPolicy
	changed bool
}

func constructResource(_ context.Context, plan *NetworkCloudFirewallPolicyResourceModel) *policyRequest {
	return constructPolicy(plan, nil)
}

func constructPolicy(plan, state *NetworkCloudFirewallPolicyResourceModel) *policyRequest {
	b := &policyRequest{CloudFirewallPolicy: graph.NewCloudFirewallPolicy()}
	if state == nil || !plan.Name.Equal(state.Name) {
		convert.FrameworkToGraphString(plan.Name, b.SetName)
		b.changed = true
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		convert.FrameworkToGraphString(plan.Description, b.SetDescription)
		if plan.Description.IsNull() {
			// Match the existing Service Principal construction pattern for explicit JSON null.
			b.SetAdditionalData(map[string]any{"description": nil})
		}
		b.changed = true
	}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := graph.NewCloudFirewallPolicySettings()
		action := graph.ALLOW_CLOUDFIREWALLACTION
		settings.SetDefaultAction(&action)
		b.SetSettings(settings)
		b.changed = true
	}
	return b
}
