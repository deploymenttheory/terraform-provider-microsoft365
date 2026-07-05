package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *OnPremisesConnectorGroupResourceModel, remoteResource *connectorGroupResponse) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.id).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.id)
	data.Name = convert.GraphToFrameworkString(remoteResource.name)
	data.ConnectorGroupType = convert.GraphToFrameworkString(remoteResource.connectorGroupType)
	data.IsDefault = convert.GraphToFrameworkBool(remoteResource.isDefault)
	data.Region = convert.GraphToFrameworkString(remoteResource.region)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
