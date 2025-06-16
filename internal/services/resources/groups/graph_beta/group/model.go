// REF: https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta
package graphBetaGroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	Description                   types.String   `tfsdk:"description"`
	MailNickname                  types.String   `tfsdk:"mail_nickname"`
	MailEnabled                   types.Bool     `tfsdk:"mail_enabled"`
	SecurityEnabled               types.Bool     `tfsdk:"security_enabled"`
	GroupTypes                    types.Set      `tfsdk:"group_types"`
	Visibility                    types.String   `tfsdk:"visibility"`
	IsAssignableToRole            types.Bool     `tfsdk:"is_assignable_to_role"`
	MembershipRule                types.String   `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String   `tfsdk:"membership_rule_processing_state"`
	CreatedDateTime               types.String   `tfsdk:"created_date_time"`
	Mail                          types.String   `tfsdk:"mail"`
	ProxyAddresses                types.Set      `tfsdk:"proxy_addresses"`
	OnPremisesSyncEnabled         types.Bool     `tfsdk:"on_premises_sync_enabled"`
	PreferredDataLocation         types.String   `tfsdk:"preferred_data_location"`
	PreferredLanguage             types.String   `tfsdk:"preferred_language"`
	Theme                         types.String   `tfsdk:"theme"`
	Classification                types.String   `tfsdk:"classification"`
	ExpirationDateTime            types.String   `tfsdk:"expiration_date_time"`
	RenewedDateTime               types.String   `tfsdk:"renewed_date_time"`
	SecurityIdentifier            types.String   `tfsdk:"security_identifier"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
} 