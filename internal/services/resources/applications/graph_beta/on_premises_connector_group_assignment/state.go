package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *OnPremisesConnectorGroupAssignmentResourceModel, remoteResource graphmodels.ConnectorGroupable) {
	tflog.Debug(ctx, "Starting to map remote connector group assignment state to Terraform state", map[string]any{
		"application_id":     data.ApplicationID.ValueString(),
		"connector_group_id": data.ConnectorGroupID.ValueString(),
	})

	data.ID = types.StringValue(compositeID(data.ApplicationID.ValueString(), data.ConnectorGroupID.ValueString()))
	data.ConnectorGroupName = convert.GraphToFrameworkString(remoteResource.GetName())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
