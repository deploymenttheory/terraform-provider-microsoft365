package graphBetaNetworkContentPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(_ context.Context, data *NetworkContentPolicyResourceModel, remote *contentPolicyResponse) {
	if remote == nil {
		return
	}
	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Description = convert.GraphToFrameworkString(remote.description)
	data.DefaultAction = convert.GraphToFrameworkString(remote.defaultAction)
}
