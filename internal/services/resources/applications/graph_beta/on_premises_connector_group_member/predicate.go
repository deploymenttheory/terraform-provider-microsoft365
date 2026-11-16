package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// connectorGroupConsistencyPredicate reports whether the read-back state shows
// the connector in the expected connector group. Graph can briefly return the
// previous group after a move.
func connectorGroupConsistencyPredicate(expectedConnectorGroupID string) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var connectorGroupID types.String
		if diags := state.GetAttribute(ctx, path.Root("connector_group_id"), &connectorGroupID); diags.HasError() {
			return false
		}

		return sameID(connectorGroupID.ValueString(), expectedConnectorGroupID)
	}
}
