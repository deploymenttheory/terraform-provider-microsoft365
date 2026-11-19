package graphBetaNetworkConditionalAccessSettings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapRemoteState(
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

func consistencyPredicate(expected types.String) func(context.Context, tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkConditionalAccessSettingsResourceModel
		if state.Get(ctx, &actual).HasError() {
			return false
		}
		return actual.ID.ValueString() == singletonID && actual.SignalingStatus.Equal(expected)
	}
}

// decodeSettingsResponse validates and decodes the raw conditional-access settings response.
func decodeSettingsResponse(result any, err error) (*conditionalAccessSettingsResponse, error) {
	if err != nil {
		return nil, err
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 {
		return nil, errInvalidSettingsResponse
	}
	var settings *conditionalAccessSettingsResponse
	if err := json.Unmarshal(raw, &settings); err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidSettingsResponse, err)
	}
	if settings == nil {
		return nil, errInvalidSettingsResponse
	}
	return settings, nil
}
