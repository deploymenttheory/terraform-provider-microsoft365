package graphBetaNetworkDefaultForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// constructResource builds the raw JSON PATCH body for the built-in profile. It carries only the configured settings
// that differ from current, so adopting a profile whose settings already match sends no request.
func constructResource(ctx context.Context, plan, current NetworkDefaultForwardingProfileResourceModel) (map[string]any, diag.Diagnostics) {
	body := map[string]any{}
	planned := sharedmodels.ForwardingProfileSettings{
		State:                    plan.State,
		ClientFallbackAction:     plan.ClientFallbackAction,
		ForwardingSourceMatching: plan.ForwardingSourceMatching,
		LocalNetworkSettings:     plan.LocalNetworkSettings,
		AssignmentRules:          plan.AssignmentRules,
	}
	prior := sharedmodels.ForwardingProfileSettings{
		State:                    current.State,
		ClientFallbackAction:     current.ClientFallbackAction,
		ForwardingSourceMatching: current.ForwardingSourceMatching,
		LocalNetworkSettings:     current.LocalNetworkSettings,
		AssignmentRules:          current.AssignmentRules,
	}
	diags := sharedmodels.SettingsPatch(ctx, body, planned, &prior)
	return body, diags
}
