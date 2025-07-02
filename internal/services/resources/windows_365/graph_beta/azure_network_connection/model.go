// REF: https://learn.microsoft.com/en-us/graph/api/virtualendpoint-list-onpremisesconnections?view=graph-rest-beta&tabs=http
package graphBetaAzureNetworkConnection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcOnPremisesConnectionResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	DisplayName        types.String   `tfsdk:"display_name"`
	ConnectionType     types.String   `tfsdk:"connection_type"`
	AdDomainName       types.String   `tfsdk:"ad_domain_name"`
	AdDomainUsername   types.String   `tfsdk:"ad_domain_username"`
	AdDomainPassword   types.String   `tfsdk:"ad_domain_password"`
	OrganizationalUnit types.String   `tfsdk:"organizational_unit"`
	ResourceGroupId    types.String   `tfsdk:"resource_group_id"`
	SubnetId           types.String   `tfsdk:"subnet_id"`
	SubscriptionId     types.String   `tfsdk:"subscription_id"`
	VirtualNetworkId   types.String   `tfsdk:"virtual_network_id"`
	HealthCheckStatus  types.String   `tfsdk:"health_check_status"`
	ManagedBy          types.String   `tfsdk:"managed_by"`
	InUse              types.Bool     `tfsdk:"in_use"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
