package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the connector and the connector group it
// currently belongs to into Terraform state.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data *OnPremisesConnectorGroupMemberResourceModel,
	connector graphmodels.Connectorable,
	connectorGroup graphmodels.ConnectorGroupable,
) {
	if connector == nil || connectorGroup == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote connector group membership to Terraform state", map[string]any{
		"connector_id": convert.GraphToFrameworkString(connector.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(connector.GetId())
	data.ConnectorID = convert.GraphToFrameworkString(connector.GetId())
	data.MachineName = convert.GraphToFrameworkString(connector.GetMachineName())
	data.ExternalIP = convert.GraphToFrameworkString(connector.GetExternalIp())
	data.ConnectorGroupID = convert.GraphToFrameworkString(connectorGroup.GetId())
	data.ConnectorGroupName = convert.GraphToFrameworkString(connectorGroup.GetName())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
