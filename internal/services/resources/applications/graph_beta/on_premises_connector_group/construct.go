package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource builds the connectorGroup create (state nil) or update body.
//
// Microsoft Learn documents connectorGroup create/update at:
// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post?view=graph-rest-beta
// https://learn.microsoft.com/en-us/graph/api/connectorgroup-update?view=graph-rest-beta
//
// Direct API verification on 2026-07-05 showed that connectorGroupType and
// isDefault are system-managed values. Graph accepted those fields in create
// payloads but ignored them, and PATCH rejected them as read-only. Keep this
// request body intentionally narrow so Terraform only sends writable fields.
//
// Region is writable according to Learn, but only while the connector group
// has no assigned connectors or applications:
// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
// An update therefore sends region only when it changed, which avoids
// unnecessary Graph validation of a field that can become non-writable after
// assignments are added.
func constructResource(ctx context.Context, plan, state *OnPremisesConnectorGroupResourceModel) connectorGroupRequest {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody connectorGroupRequest

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		requestBody.Name = plan.Name.ValueStringPointer()
	}

	regionChanged := state == nil || state.Region.IsNull() || state.Region.IsUnknown() || plan.Region.ValueString() != state.Region.ValueString()
	if !plan.Region.IsNull() && !plan.Region.IsUnknown() && regionChanged {
		requestBody.Region = plan.Region.ValueStringPointer()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName), map[string]any{"request_body": requestBody})

	return requestBody
}
