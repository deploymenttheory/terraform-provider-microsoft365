package graphBetaNetworkExplicitForwardProxy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkExplicitForwardProxyResourceModel holds the managed configuration.
type NetworkExplicitForwardProxyResourceModel struct {
	ID                            types.String                 `tfsdk:"id"`
	InternetAccess                *InternetAccessResourceModel `tfsdk:"internet_access"`
	ProxyAutoConfigurationFileURL types.String                 `tfsdk:"proxy_auto_configuration_file_url"`
	Timeouts                      timeouts.Value               `tfsdk:"timeouts"`
}

// InternetAccessResourceModel follows the internetAccess API object.
type InternetAccessResourceModel struct {
	IsEnabled                        types.Bool   `tfsdk:"is_enabled"`
	IsSourceIPSessionAffinityEnabled types.Bool   `tfsdk:"is_source_ip_session_affinity_enabled"`
	SourceIPSessionAffinityOptions   types.String `tfsdk:"source_ip_session_affinity_options"`
}

type proxyResponse struct {
	ProxyAutoConfigurationFileURL *string `json:"proxyAutoConfigurationFileUrl"`
	InternetAccess                *struct {
		IsEnabled                        *bool   `json:"isEnabled"`
		IsSourceIPSessionAffinityEnabled *bool   `json:"isSourceIpSessionAffinityEnabled"`
		SourceIPSessionAffinityOptions   *string `json:"sourceIpSessionAffinityOptions"`
	} `json:"internetAccess"`
}
