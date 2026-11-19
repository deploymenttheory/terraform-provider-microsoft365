// REF: https://learn.microsoft.com/en-us/graph/api/resources/adminconsentrequestpolicy?view=graph-rest-beta
package graphBetaAdminConsentRequestPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AdminConsentRequestPolicyResourceModel represents the schema for the admin consent request policy resource.
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

// ReviewerModel represents a reviewer query for the admin consent request policy.
type ReviewerModel struct {
	Query     types.String `tfsdk:"query"`
	QueryType types.String `tfsdk:"query_type"`
	QueryRoot types.String `tfsdk:"query_root"`
}

var reviewerAttributeTypes = map[string]attr.Type{"query": types.StringType, "query_type": types.StringType, "query_root": types.StringType}
