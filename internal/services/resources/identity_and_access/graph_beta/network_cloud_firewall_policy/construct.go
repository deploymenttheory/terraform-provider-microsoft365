package graphBetaNetworkCloudFirewallPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource builds the cloud firewall policy body. Create passes a nil state and sends every field; Update passes the
// prior state and sends only the fields that changed.
func constructResource(_ context.Context, plan, state *NetworkCloudFirewallPolicyResourceModel) graph.CloudFirewallPolicyable {
	body := graph.NewCloudFirewallPolicy()
	if state == nil || !plan.Name.Equal(state.Name) {
		convert.FrameworkToGraphString(plan.Name, body.SetName)
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		convert.FrameworkToGraphString(plan.Description, body.SetDescription)
		if plan.Description.IsNull() {
			// The generated setter omits nil strings, so an explicit JSON null goes through additional data.
			body.SetAdditionalData(map[string]any{"description": nil})
		}
	}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := graph.NewCloudFirewallPolicySettings()
		action := graph.ALLOW_CLOUDFIREWALLACTION
		settings.SetDefaultAction(&action)
		body.SetSettings(settings)
	}
	return body
}
