// https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-policy?view=graph-rest-beta

package graphBetaWindowsUpdatesAutopatchPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchPolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	ApprovalRules        types.Set      `tfsdk:"approval_rules"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type ApprovalRuleModel struct {
	DeferralInDays types.Int32  `tfsdk:"deferral_in_days"`
	Classification types.String `tfsdk:"classification"`
	Cadence        types.String `tfsdk:"cadence"`
}
