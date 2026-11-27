// REF: https://learn.microsoft.com/en-us/graph/api/tenantrelationship-findtenantinformationbytenantid?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/tenantrelationship-findtenantinformationbydomainname?view=graph-rest-beta

package graphBetaTenantInformation

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantInformationDataSourceModel struct {
	ID         types.String            `tfsdk:"id"`
	TenantID   types.String            `tfsdk:"tenant_id"`
	DomainName types.String            `tfsdk:"domain_name"`
	Items      []TenantInformationItem `tfsdk:"items"`
	Timeouts   timeouts.Value          `tfsdk:"timeouts"`
}
type TenantInformationItem struct {
	TenantID            types.String `tfsdk:"tenant_id"`
	DisplayName         types.String `tfsdk:"display_name"`
	DefaultDomainName   types.String `tfsdk:"default_domain_name"`
	FederationBrandName types.String `tfsdk:"federation_brand_name"`
}
