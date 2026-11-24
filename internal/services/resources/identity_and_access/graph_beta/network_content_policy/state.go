package graphBetaNetworkContentPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(_ context.Context, data *NetworkContentPolicyResourceModel, remote *contentPolicyResponse) {
	if remote == nil {
		return
	}
	var defaultAction *string
	if remote.Settings != nil {
		defaultAction = remote.Settings.DefaultAction
	}
	data.ID = convert.GraphToFrameworkString(remote.ID)
	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.Description = convert.GraphToFrameworkString(remote.Description)
	data.DefaultAction = convert.GraphToFrameworkString(defaultAction)
}
