package graphBetaNetworkDeviceForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// mapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
func mapRemoteStateToTerraform(
	ctx context.Context,
	state *NetworkDeviceForwardingProfileResourceModel,
	remote *sharedmodels.ForwardingProfileWire,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if remote.IsCustomProfile != nil && !*remote.IsCustomProfile {
		diags.AddError(
			"Not a custom forwarding profile",
			"This profile is a built-in pre-logon profile and cannot be managed by this resource.",
		)
		return diags
	}
	state.ID = types.StringPointerValue(remote.ID)
	state.Name = types.StringPointerValue(remote.Name)
	state.Description = types.StringPointerValue(remote.Description)
	state.TrafficForwardingType = types.StringPointerValue(remote.TrafficForwardingType)
	state.Priority = types.Int32PointerValue(remote.Priority)
	state.Version = types.StringPointerValue(remote.Version)
	settings, d := sharedmodels.SettingsFromWire(ctx, remote)
	diags.Append(d...)
	state.State = settings.State
	state.ClientFallbackAction = settings.ClientFallbackAction
	state.LocalNetworkSettings = settings.LocalNetworkSettings
	state.AssignmentRules = settings.AssignmentRules
	return diags
}
