package graphBetaNetworkConditionalAccessSettings

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the signaling status onto state, or reports an invalid response and leaves state unchanged.
func MapRemoteStateToTerraform(
	data *NetworkConditionalAccessSettingsResourceModel,
	remote *conditionalAccessSettingsResponse,
	diagnostics *diag.Diagnostics,
) {
	if remote == nil || remote.SignalingStatus == nil {
		diagnostics.AddError(
			"Invalid conditional access settings response",
			"Graph returned no signalingStatus. Existing Terraform state has been preserved.",
		)
		return
	}
	status := convert.GraphToFrameworkString(remote.SignalingStatus)
	if status.ValueString() != "enabled" && status.ValueString() != "disabled" {
		diagnostics.AddError(
			"Invalid conditional access settings response",
			"Graph returned an unsupported signalingStatus. Existing Terraform state has been preserved.",
		)
		return
	}
	data.ID = types.StringValue(singletonID)
	data.SignalingStatus = status
}
