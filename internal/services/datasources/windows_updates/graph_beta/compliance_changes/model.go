package graphBetaWindowsUpdatesComplianceChanges

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ComplianceChangesDataSourceModel struct {
	UpdatePolicyId    types.String        `tfsdk:"update_policy_id"`
	ComplianceChanges []ComplianceChange  `tfsdk:"compliance_changes"`
	Timeouts          timeouts.Value      `tfsdk:"timeouts"`
}

type ComplianceChange struct {
	ID              types.String         `tfsdk:"id"`
	CreatedDateTime types.String         `tfsdk:"created_date_time"`
	IsRevoked       types.Bool           `tfsdk:"is_revoked"`
	RevokedDateTime types.String         `tfsdk:"revoked_date_time"`
	Content         *ComplianceContent   `tfsdk:"content"`
	DeploymentSettings *DeploymentSettings `tfsdk:"deployment_settings"`
}

type ComplianceContent struct {
	CatalogEntryId   types.String `tfsdk:"catalog_entry_id"`
	CatalogEntryType types.String `tfsdk:"catalog_entry_type"`
}

type DeploymentSettings struct {
	Schedule *ScheduleSettings `tfsdk:"schedule"`
}

type ScheduleSettings struct {
	StartDateTime  types.String    `tfsdk:"start_date_time"`
	GradualRollout *GradualRollout `tfsdk:"gradual_rollout"`
}

type GradualRollout struct {
	DurationBetweenOffers types.String `tfsdk:"duration_between_offers"`
	DevicesPerOffer       types.Int32  `tfsdk:"devices_per_offer"`
}
