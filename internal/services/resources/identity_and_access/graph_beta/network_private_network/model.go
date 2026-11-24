package graphBetaNetworkPrivateNetwork

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkPrivateNetworkResourceModel represents the Terraform schema for
// Microsoft Entra Global Secure Access private networks.
type NetworkPrivateNetworkResourceModel struct {
	ID                          types.String                      `tfsdk:"id"`
	Name                        types.String                      `tfsdk:"name"`
	DNSResolutionIdentification *DNSResolutionIdentificationModel `tfsdk:"dns_resolution_identification"`
	AppIDs                      types.Set                         `tfsdk:"app_ids"`
	Timeouts                    timeouts.Value                    `tfsdk:"timeouts"`
}

type DNSResolutionIdentificationModel struct {
	DNSServers            types.Set    `tfsdk:"dns_servers"`
	FQDNToResolve         types.String `tfsdk:"fqdn_to_resolve"`
	ExpectedIPResolutions types.Set    `tfsdk:"expected_ip_resolutions"`
}

type ExpectedIPResolutionModel struct {
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
	BeginAddress types.String `tfsdk:"begin_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

// privateNetworkRequest is the portal-observed privateNetwork create and update body. @odata.type is sent on
// create only. AppIDs is omitted when app_ids is null and sent as [] when it is empty.
type privateNetworkRequest struct {
	ODataType              string                               `json:"@odata.type,omitempty"` //nolint:tagliatelle // OData annotation name.
	Name                   *string                              `json:"name,omitempty"`
	AppIDs                 *[]string                            `json:"appIds,omitempty"`
	NetworkIdentifications []dnsResolutionIdentificationRequest `json:"networkIdentifications"`
}

// dnsResolutionIdentificationRequest is the networkaccess.dnsResolutionIdentification network identification.
type dnsResolutionIdentificationRequest struct {
	ODataType             string                `json:"@odata.type"` //nolint:tagliatelle // OData annotation name.
	ServerAddresses       []ipResolutionRequest `json:"serverAddresses"`
	FQDNToResolve         *ipResolutionRequest  `json:"fqdnToResolve"`
	ExpectedIPResolutions []ipResolutionRequest `json:"expectedIpResolutions"`
}

// ipResolutionRequest is a networkaccess.ipAddress, ipSubnet, ipRange or fqdn value, discriminated by @odata.type.
type ipResolutionRequest struct {
	ODataType    string  `json:"@odata.type"` //nolint:tagliatelle // OData annotation name.
	Value        *string `json:"value,omitempty"`
	BeginAddress *string `json:"beginAddress,omitempty"`
	EndAddress   *string `json:"endAddress,omitempty"`
}

// privateNetworkResponse holds the privateNetwork fields this resource reads.
type privateNetworkResponse struct {
	ID                     *string                                `json:"id"`
	Name                   *string                                `json:"name"`
	AppIDs                 []string                               `json:"appIds"`
	NetworkIdentifications []*dnsResolutionIdentificationResponse `json:"networkIdentifications"`
}

// dnsResolutionIdentificationResponse is the networkaccess.dnsResolutionIdentification network identification.
type dnsResolutionIdentificationResponse struct {
	ServerAddresses       []ipResolutionResponse `json:"serverAddresses"`
	FQDNToResolve         *ipResolutionResponse  `json:"fqdnToResolve"`
	ExpectedIPResolutions []ipResolutionResponse `json:"expectedIpResolutions"`
}

// ipResolutionResponse is a networkaccess.ipAddress, ipSubnet, ipRange or fqdn value, discriminated by @odata.type.
type ipResolutionResponse struct {
	ODataType    *string `json:"@odata.type"` //nolint:tagliatelle // OData annotation name.
	Value        *string `json:"value"`
	BeginAddress *string `json:"beginAddress"`
	EndAddress   *string `json:"endAddress"`
}
