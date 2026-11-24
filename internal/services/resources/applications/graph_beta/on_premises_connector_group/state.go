package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *OnPremisesConnectorGroupResourceModel, remoteResource *connectorGroupResponse) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.ID).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.ID)
	data.Name = convert.GraphToFrameworkString(remoteResource.Name)
	data.ConnectorGroupType = convert.GraphToFrameworkString(remoteResource.ConnectorGroupType)
	data.IsDefault = convert.GraphToFrameworkBool(remoteResource.IsDefault)
	data.Region = convert.GraphToFrameworkString(remoteResource.Region)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
