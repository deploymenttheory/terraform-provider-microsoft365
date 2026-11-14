// REF: https://learn.microsoft.com/en-us/graph/api/resources/connectedorganization?view=graph-rest-beta
package graphBetaIdentityAndAccessConnectedOrganization

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConnectedOrganizationResourceModel represents the schema for the ConnectedOrganization resource
type ConnectedOrganizationResourceModel struct {
	ID               types.String          `tfsdk:"id"`
	DisplayName      types.String          `tfsdk:"display_name"`
	Description      types.String          `tfsdk:"description"`
	State            types.String          `tfsdk:"state"`
	IdentitySources  []IdentitySourceModel `tfsdk:"identity_sources"`
	CreatedBy        types.String          `tfsdk:"created_by"`
	CreatedDateTime  types.String          `tfsdk:"created_date_time"`
	ModifiedBy       types.String          `tfsdk:"modified_by"`
	ModifiedDateTime types.String          `tfsdk:"modified_date_time"`
	Timeouts         timeouts.Value        `tfsdk:"timeouts"`
}

// IdentitySourceModel represents one polymorphic identitySource of a connected organization.
// Which of domain_name, tenant_id and issuer_uri apply depends on odata_type.
type IdentitySourceModel struct {
	ODataType   types.String `tfsdk:"odata_type"`
	DisplayName types.String `tfsdk:"display_name"`
	DomainName  types.String `tfsdk:"domain_name"`
	TenantID    types.String `tfsdk:"tenant_id"`
	IssuerURI   types.String `tfsdk:"issuer_uri"`
}
