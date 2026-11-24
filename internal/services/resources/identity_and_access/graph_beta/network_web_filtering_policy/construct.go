package graphBetaNetworkWebFilteringPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultActionAllowODataType = "#microsoft.graph.networkaccess.webFilteringActionAllow"
	defaultActionBlockODataType = "#microsoft.graph.networkaccess.webFilteringActionBlock"
)

// constructResource builds the portal-observed web filtering policy
// payload for /networkaccess/webFilteringPolicies.
//
// Microsoft Graph beta currently documents the older/generic
// networkaccess.filteringPolicy shape, not this webFilteringPolicy surface:
// https://learn.microsoft.com/graph/api/resources/networkaccess-filteringpolicy
//
// The Entra Global Secure Access Web content filtering blade sends
// settings.defaultAction.@odata.type plus an empty policyRules array on create.
// policyRules is intentionally create-only here because subsequent rule
// management happens through the child /policyRules resource.
func constructResource(ctx context.Context, data *NetworkWebFilteringPolicyResourceModel, includePolicyRules bool) webFilteringPolicy {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	defaultAction := defaultActionAllowODataType
	if data.DefaultAction.ValueString() == "block" {
		defaultAction = defaultActionBlockODataType
	}

	requestBody := webFilteringPolicy{
		Settings: &webFilteringPolicySettings{DefaultAction: &webFilteringDefaultAction{ODataType: defaultAction}},
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody.Name = data.Name.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody.Description = data.Description.ValueStringPointer()
	}
	if includePolicyRules {
		requestBody.PolicyRules = &[]struct{}{}
	}

	return requestBody
}
