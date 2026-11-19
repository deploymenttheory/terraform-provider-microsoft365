package graphBetaNetworkProxyAutoConfiguration

import (
	"fmt"

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

// decodeCreateResponse validates and maps the Microsoft Graph response.
func decodeCreateResponse(response pacResponse, err error) (*pacResponse, error) {
	if err != nil {
		return nil, fmt.Errorf("create custom PAC: %w", err)
	}
	if response.ID == nil || *response.ID == "" {
		return nil, errInvalidPACResponse
	}
	return &response, nil
}
