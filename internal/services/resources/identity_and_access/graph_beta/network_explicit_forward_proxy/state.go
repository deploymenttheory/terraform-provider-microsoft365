package graphBetaNetworkExplicitForwardProxy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapResponse(state *NetworkExplicitForwardProxyResourceModel, remote *proxyResponse) error {
	if remote == nil || remote.InternetAccess == nil {
		return errInvalidProxyResponse
	}
	access := remote.InternetAccess
	if access.IsEnabled == nil || access.IsSourceIPSessionAffinityEnabled == nil ||
		access.SourceIPSessionAffinityOptions == nil {
		return errInvalidProxyResponse
	}
	options := *access.SourceIPSessionAffinityOptions
	if options == "useHttpHeader,useSessionId" {
		options = "useSessionId,useHttpHeader"
	}
	if !validAffinity(*access.IsSourceIPSessionAffinityEnabled, options) {
		return errInvalidProxyResponse
	}
	state.ID = types.StringValue("explicitForwardProxyConfig")
	state.InternetAccess = &InternetAccessResourceModel{
		IsEnabled: convert.GraphToFrameworkBool(access.IsEnabled),
		IsSourceIPSessionAffinityEnabled: convert.GraphToFrameworkBool(
			access.IsSourceIPSessionAffinityEnabled,
		),
		SourceIPSessionAffinityOptions: types.StringValue(options),
	}
	state.ProxyAutoConfigurationFileURL = convert.GraphToFrameworkString(
		remote.ProxyAutoConfigurationFileURL,
	)
	return nil
}
