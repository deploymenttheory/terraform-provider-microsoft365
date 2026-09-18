package graphBetaNetworkProxyAutoConfiguration

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapResponse(state *NetworkProxyAutoConfigurationResourceModel, remote *pacResponse) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil ||
		remote.Content == nil ||
		remote.IsEnabled == nil {
		return errInvalidPACResponse
	}
	state.ID = convert.GraphToFrameworkString(remote.ID)
	state.Name = convert.GraphToFrameworkString(remote.Name)
	state.Content = convert.GraphToFrameworkString(remote.Content)
	state.IsEnabled = convert.GraphToFrameworkBool(remote.IsEnabled)
	state.CreatedDateTime = convert.GraphToFrameworkString(remote.CreatedDateTime)
	state.LastModifiedDateTime = convert.GraphToFrameworkString(remote.LastModifiedDateTime)
	return nil
}
