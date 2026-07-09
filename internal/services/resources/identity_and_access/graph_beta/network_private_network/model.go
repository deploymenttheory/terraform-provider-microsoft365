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
