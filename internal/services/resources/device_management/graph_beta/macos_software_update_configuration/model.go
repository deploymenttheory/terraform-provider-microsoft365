// ref: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdateconfiguration?view=graph-rest-beta
package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MacOSSoftwareUpdateConfigurationResourceModel struct {
	ID                                 types.String   `tfsdk:"id"`
	DisplayName                        types.String   `tfsdk:"display_name"`
	Description                        types.String   `tfsdk:"description"`
	RoleScopeTagIds                    types.Set      `tfsdk:"role_scope_tag_ids"`
	UpdateScheduleType                 types.String   `tfsdk:"update_schedule_type"`
	CriticalUpdateBehavior             types.String   `tfsdk:"critical_update_behavior"`
	ConfigDataUpdateBehavior           types.String   `tfsdk:"config_data_update_behavior"`
	FirmwareUpdateBehavior             types.String   `tfsdk:"firmware_update_behavior"`
	AllOtherUpdateBehavior             types.String   `tfsdk:"all_other_update_behavior"`
	UpdateTimeWindowUtcOffsetInMinutes types.Int32    `tfsdk:"update_time_window_utc_offset_in_minutes"`
	CustomUpdateTimeWindows            types.List     `tfsdk:"custom_update_time_windows"`
	MaxUserDeferralsCount              types.Int32    `tfsdk:"max_user_deferrals_count"`
	Priority                           types.String   `tfsdk:"priority"`
	Assignments                        types.Set      `tfsdk:"assignments"`
	Timeouts                           timeouts.Value `tfsdk:"timeouts"`
}
