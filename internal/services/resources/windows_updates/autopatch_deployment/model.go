package graphBetaWindowsUpdatesAutopatchDeployment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchDeploymentResourceModel struct {
	ID                   types.String        `tfsdk:"id"`
	Content              *DeploymentContent  `tfsdk:"content"`
	Settings             *DeploymentSettings `tfsdk:"settings"`
	CreatedDateTime      types.String        `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String        `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value      `tfsdk:"timeouts"`
}

type DeploymentContent struct {
	CatalogEntryId   types.String `tfsdk:"catalog_entry_id"`
	CatalogEntryType types.String `tfsdk:"catalog_entry_type"`
}

type DeploymentSettings struct {
	Schedule             *ScheduleSettings             `tfsdk:"schedule"`
	Monitoring           *MonitoringSettings           `tfsdk:"monitoring"`
	UserExperience       *UserExperienceSettings       `tfsdk:"user_experience"`
	Expedite             *ExpediteSettings             `tfsdk:"expedite"`
	ContentApplicability *ContentApplicabilitySettings `tfsdk:"content_applicability"`
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
	MonitoringRules types.Set `tfsdk:"monitoring_rules"`
}

type MonitoringRule struct {
	Signal    types.String `tfsdk:"signal"`
	Threshold types.Int32  `tfsdk:"threshold"`
	Action    types.String `tfsdk:"action"`
}

type UserExperienceSettings struct {
	DaysUntilForcedReboot types.Int32 `tfsdk:"days_until_forced_reboot"`
	OfferAsOptional       types.Bool  `tfsdk:"offer_as_optional"`
}

type ExpediteSettings struct {
	IsExpedited     types.Bool `tfsdk:"is_expedited"`
	IsReadinessTest types.Bool `tfsdk:"is_readiness_test"`
}

type ContentApplicabilitySettings struct {
	Safeguard *SafeguardSettings `tfsdk:"safeguard"`
}

type SafeguardSettings struct {
	DisabledSafeguardProfiles types.Set `tfsdk:"disabled_safeguard_profiles"`
}

type SafeguardProfile struct {
	Category types.String `tfsdk:"category"`
}
