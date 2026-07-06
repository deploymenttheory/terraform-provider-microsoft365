package graphBetaNetworkWebFilteringPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkWebFilteringPolicyResourceModel, remoteResource *webFilteringPolicyResponse) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.id)
	data.Name = convert.GraphToFrameworkString(remoteResource.name)
	data.Description = convert.GraphToFrameworkString(remoteResource.description)
	data.DefaultAction = convert.GraphToFrameworkString(terraformDefaultAction(remoteResource.defaultAction))

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
