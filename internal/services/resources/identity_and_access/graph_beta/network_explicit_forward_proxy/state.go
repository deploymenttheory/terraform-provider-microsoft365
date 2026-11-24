package graphBetaNetworkExplicitForwardProxy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the explicit forward proxy configuration to Terraform state. It returns
// errInvalidProxyResponse, leaving state unchanged, when Graph omits a required setting or returns unsupported
// session affinity options.
func MapRemoteStateToTerraform(state *NetworkExplicitForwardProxyResourceModel, remote *proxyResponse) error {
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
	if wantEnabled, ok := affinityOptionsEnabled[options]; !ok || wantEnabled != *access.IsSourceIPSessionAffinityEnabled {
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
