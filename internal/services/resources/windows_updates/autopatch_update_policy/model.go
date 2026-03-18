// https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-updatepolicy?view=graph-rest-beta

package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchUpdatePolicyResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	CreatedDateTime       types.String   `tfsdk:"created_date_time"`
	AudienceId            types.String   `tfsdk:"audience_id"`
	ComplianceChanges     types.Bool     `tfsdk:"compliance_changes"`
	ComplianceChangeRules types.Set      `tfsdk:"compliance_change_rules"`
	DeploymentSettings    types.Object   `tfsdk:"deployment_settings"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

type ComplianceChangeRuleModel struct {
	ContentFilter                *ContentFilterModel `tfsdk:"content_filter"`
	DurationBeforeDeploymentStart types.String        `tfsdk:"duration_before_deployment_start"`
	CreatedDateTime               types.String        `tfsdk:"created_date_time"`
	LastEvaluatedDateTime         types.String        `tfsdk:"last_evaluated_date_time"`
	LastModifiedDateTime          types.String        `tfsdk:"last_modified_date_time"`
}

type ContentFilterModel struct {
	FilterType types.String `tfsdk:"filter_type"`
}

type DeploymentSettingsModel struct {
	Schedule types.Object `tfsdk:"schedule"`
}

type ScheduleSettingsModel struct {
	StartDateTime  types.String `tfsdk:"start_date_time"`
	GradualRollout types.Object `tfsdk:"gradual_rollout"`
}

type GradualRolloutModel struct {
	DurationBetweenOffers types.String `tfsdk:"duration_between_offers"`
	DevicesPerOffer       types.Int32  `tfsdk:"devices_per_offer"`
}
