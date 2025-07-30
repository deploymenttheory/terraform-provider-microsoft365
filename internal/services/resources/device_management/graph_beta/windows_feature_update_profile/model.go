// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofile?view=graph-rest-beta
package graphBetaWindowsFeatureUpdateProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsFeatureUpdateProfileResourceModel struct {
	ID                                                types.String          `tfsdk:"id"`
	DisplayName                                       types.String          `tfsdk:"display_name"`
	Description                                       types.String          `tfsdk:"description"`
	FeatureUpdateVersion                              types.String          `tfsdk:"feature_update_version"`
	RolloutSettings                                   *RolloutSettingsModel `tfsdk:"rollout_settings"`
	CreatedDateTime                                   types.String          `tfsdk:"created_date_time"`
	LastModifiedDateTime                              types.String          `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds                                   types.Set             `tfsdk:"role_scope_tag_ids"`
	DeployableContentDisplayName                      types.String          `tfsdk:"deployable_content_display_name"`
	EndOfSupportDate                                  types.String          `tfsdk:"end_of_support_date"`
	InstallLatestWindows10OnWindows11IneligibleDevice types.Bool            `tfsdk:"install_latest_windows10_on_windows11_ineligible_device"`
	InstallFeatureUpdatesOptional                     types.Bool            `tfsdk:"install_feature_updates_optional"`
	Assignments                                       types.Set             `tfsdk:"assignments"`
	Timeouts                                          timeouts.Value        `tfsdk:"timeouts"`
}

type RolloutSettingsModel struct {
	OfferStartDateTimeInUTC types.String `tfsdk:"offer_start_date_time_in_utc"`
	OfferEndDateTimeInUTC   types.String `tfsdk:"offer_end_date_time_in_utc"`
	OfferIntervalInDays     types.Int32  `tfsdk:"offer_interval_in_days"`
}
