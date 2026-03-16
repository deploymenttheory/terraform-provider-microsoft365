package graphBetaWindowsAutopatchDeployment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsAutopatchDeploymentResourceModel struct {
	ID                   types.String        `tfsdk:"id"`
	Content              *DeploymentContent  `tfsdk:"content"`
	Settings             *DeploymentSettings `tfsdk:"settings"`
	State                types.Object        `tfsdk:"state"`
	CreatedDateTime      types.String        `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String        `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
}

type DeploymentContent struct {
	CatalogEntryId   types.String `tfsdk:"catalog_entry_id"`
	CatalogEntryType types.String `tfsdk:"catalog_entry_type"`
}

type DeploymentSettings struct {
	Schedule   *ScheduleSettings   `tfsdk:"schedule"`
	Monitoring *MonitoringSettings `tfsdk:"monitoring"`
}

type ScheduleSettings struct {
	StartDateTime  types.String    `tfsdk:"start_date_time"`
	GradualRollout *GradualRollout `tfsdk:"gradual_rollout"`
}

type GradualRollout struct {
	DurationBetweenOffers types.String `tfsdk:"duration_between_offers"`
	DevicesPerOffer       types.Int32  `tfsdk:"devices_per_offer"`
	EndDateTime           types.String `tfsdk:"end_date_time"`
}

type MonitoringSettings struct {
	MonitoringRules []MonitoringRule `tfsdk:"monitoring_rules"`
}

type MonitoringRule struct {
	Signal    types.String `tfsdk:"signal"`
	Threshold types.Int32  `tfsdk:"threshold"`
	Action    types.String `tfsdk:"action"`
}

type DeploymentState struct {
	RequestedValue types.String `tfsdk:"requested_value"`
	EffectiveValue types.String `tfsdk:"effective_value"`
}
