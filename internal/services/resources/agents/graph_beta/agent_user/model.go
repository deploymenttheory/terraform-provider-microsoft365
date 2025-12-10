// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentuser?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-update?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-delete?view=graph-rest-beta
package graphBetaAgentUser

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentUserResourceModel represents the Terraform resource model for an agent user.
type AgentUserResourceModel struct {
	// Required fields
	ID                types.String `tfsdk:"id"`
	DisplayName       types.String `tfsdk:"display_name"`
	AgentIdentityId   types.String `tfsdk:"agent_identity_id"`
	AccountEnabled    types.Bool   `tfsdk:"account_enabled"`
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
	MailNickname      types.String `tfsdk:"mail_nickname"`

	// Computed fields (read-only)
	Mail            types.String `tfsdk:"mail"`
	UserType        types.String `tfsdk:"user_type"`
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	CreationType    types.String `tfsdk:"creation_type"`

	// Optional name fields
	GivenName types.String `tfsdk:"given_name"`
	Surname   types.String `tfsdk:"surname"`

	// Optional organizational fields
	JobTitle       types.String `tfsdk:"job_title"`
	Department     types.String `tfsdk:"department"`
	CompanyName    types.String `tfsdk:"company_name"`
	OfficeLocation types.String `tfsdk:"office_location"`

	// Optional address fields
	City          types.String `tfsdk:"city"`
	State         types.String `tfsdk:"state"`
	Country       types.String `tfsdk:"country"`
	PostalCode    types.String `tfsdk:"postal_code"`
	StreetAddress types.String `tfsdk:"street_address"`

	// Optional locale fields
	UsageLocation     types.String `tfsdk:"usage_location"`
	PreferredLanguage types.String `tfsdk:"preferred_language"`

	// Relationships
	SponsorIds types.Set `tfsdk:"sponsor_ids"`

	// Terraform-specific
	HardDelete types.Bool     `tfsdk:"hard_delete"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}
