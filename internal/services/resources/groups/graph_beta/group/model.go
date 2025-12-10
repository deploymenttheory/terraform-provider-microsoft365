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
	GroupOwners                   types.Set      `tfsdk:"group_owners"`
	GroupMembers                  types.Set      `tfsdk:"group_members"`
	HardDelete                    types.Bool     `tfsdk:"hard_delete"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}
