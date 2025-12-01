package utilityMicrosoft365EndpointReference

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Microsoft365EndpointReferenceDataSourceModel describes the Terraform datasource data model
type Microsoft365EndpointReferenceDataSourceModel struct {
	Id           types.String   `tfsdk:"id"`
	Instance     types.String   `tfsdk:"instance"`
	ServiceAreas types.Set      `tfsdk:"service_areas"`
	Categories   types.Set      `tfsdk:"categories"`
	RequiredOnly types.Bool     `tfsdk:"required_only"`
	ExpressRoute types.Bool     `tfsdk:"express_route"`
	Endpoints    types.List     `tfsdk:"endpoints"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

// EndpointModel represents a single Microsoft 365 endpoint set
type EndpointModel struct {
	Id                     types.Int64  `tfsdk:"id"`
	ServiceArea            types.String `tfsdk:"service_area"`
	ServiceAreaDisplayName types.String `tfsdk:"service_area_display_name"`
	Urls                   types.List   `tfsdk:"urls"`
	IPs                    types.List   `tfsdk:"ips"`
	TCPPorts               types.String `tfsdk:"tcp_ports"`
	UDPPorts               types.String `tfsdk:"udp_ports"`
	ExpressRoute           types.Bool   `tfsdk:"express_route"`
	Category               types.String `tfsdk:"category"`
	Required               types.Bool   `tfsdk:"required"`
	Notes                  types.String `tfsdk:"notes"`
}
