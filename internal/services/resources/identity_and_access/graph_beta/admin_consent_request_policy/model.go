package graphBetaAdminConsentRequestPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AdminConsentRequestPolicyResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	IsEnabled             types.Bool     `tfsdk:"is_enabled"`
	NotifyReviewers       types.Bool     `tfsdk:"notify_reviewers"`
	RemindersEnabled      types.Bool     `tfsdk:"reminders_enabled"`
	RequestDurationInDays types.Int32    `tfsdk:"request_duration_in_days"`
	Reviewers             types.Set      `tfsdk:"reviewers"`
	Version               types.Int32    `tfsdk:"version"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
type ReviewerModel struct {
	Query     types.String `tfsdk:"query"`
	QueryType types.String `tfsdk:"query_type"`
	QueryRoot types.String `tfsdk:"query_root"`
}

var reviewerAttributeTypes = map[string]attr.Type{"query": types.StringType, "query_type": types.StringType, "query_root": types.StringType}
