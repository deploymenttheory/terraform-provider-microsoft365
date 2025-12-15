// REF: https://learn.microsoft.com/en-us/graph/api/tenantrelationship-findtenantinformationbytenantid?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/tenantrelationship-findtenantinformationbydomainname?view=graph-rest-beta

package graphBetaTenantInformation

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantInformationDataSourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	FilterType          types.String   `tfsdk:"filter_type"`
	FilterValue         types.String   `tfsdk:"filter_value"`
	TenantID            types.String   `tfsdk:"tenant_id"`
	DisplayName         types.String   `tfsdk:"display_name"`
	DefaultDomainName   types.String   `tfsdk:"default_domain_name"`
	FederationBrandName types.String   `tfsdk:"federation_brand_name"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

