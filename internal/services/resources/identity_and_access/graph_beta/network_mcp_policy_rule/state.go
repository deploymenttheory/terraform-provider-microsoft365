package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteStateToTerraform validates the response before replacing the last known state.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *NetworkMCPPolicyRuleResourceModel,
	remote *mcpPolicyRuleResponse,
) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil ||
		remote.Action == nil ||
		remote.Priority == nil ||
		remote.Settings == nil ||
		remote.Settings.Status == nil ||
		remote.ODataType == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if *remote.ODataType != "#microsoft.graph.networkaccess.mcpPolicyRule" {
		return fmt.Errorf("%w: unsupported rule type %q", errInvalidResponse, *remote.ODataType)
	}
	if *remote.Priority < 100 || *remote.Priority > math.MaxInt32 {
		return fmt.Errorf("%w: invalid priority", errInvalidResponse)
	}
	if *remote.Action != "allow" && *remote.Action != "block" {
		return fmt.Errorf("%w: unsupported action %q", errInvalidResponse, *remote.Action)
	}
	status := *remote.Settings.Status
	if status != "enabled" && status != "disabled" {
		return fmt.Errorf(
			"%w: unsupported status %q; enabled cannot be determined",
			errInvalidResponse,
			status,
		)
	}
	if len(remote.Conditions) == 0 {
		return fmt.Errorf("%w: missing matchingConditions", errInvalidResponse)
	}
	conditions, err := conditionsToState(ctx, remote.Conditions)
	if err != nil {
		return err
	}
	data.ID = types.StringValue(*remote.ID)
	data.Name = types.StringValue(*remote.Name)
	data.Description = types.StringPointerValue(remote.Description)
	data.Action = types.StringValue(*remote.Action)
	data.Priority = types.Int32Value(int32(*remote.Priority))
	data.Status = types.StringValue(status)
	data.Enabled = types.BoolValue(status == "enabled")
	data.MatchingConditions = conditions
	return nil
}
