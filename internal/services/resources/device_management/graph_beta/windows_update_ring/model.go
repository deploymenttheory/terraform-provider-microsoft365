// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateforbusinessconfiguration?view=graph-rest-beta
package graphBetaWindowsUpdateRing

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsUpdateRingResourceModel defines the model for Windows Update Ring resource
type WindowsUpdateRingResourceModel struct {
	ID                                                      types.String            `tfsdk:"id"`
	DisplayName                                             types.String            `tfsdk:"display_name"`
	Description                                             types.String            `tfsdk:"description"`
	RoleScopeTagIds                                         types.Set               `tfsdk:"role_scope_tag_ids"`
	MicrosoftUpdateServiceAllowed                           types.Bool              `tfsdk:"microsoft_update_service_allowed"`
	DriversExcluded                                         types.Bool              `tfsdk:"drivers_excluded"`
	QualityUpdatesDeferralPeriodInDays                      types.Int32             `tfsdk:"quality_updates_deferral_period_in_days"`
	FeatureUpdatesDeferralPeriodInDays                      types.Int32             `tfsdk:"feature_updates_deferral_period_in_days"`
	FeatureUpdatesPauseExpiryDateTime                       types.String            `tfsdk:"feature_updates_pause_expiry_date_time"`
	FeatureUpdatesPauseStartDate                            types.String            `tfsdk:"feature_updates_pause_start_date"`
	FeatureUpdatesRollbackStartDateTime                     types.String            `tfsdk:"feature_updates_rollback_start_date_time"`
	QualityUpdatesPauseExpiryDateTime                       types.String            `tfsdk:"quality_updates_pause_expiry_date_time"`
	QualityUpdatesPauseStartDate                            types.String            `tfsdk:"quality_updates_pause_start_date"`
	QualityUpdatesRollbackStartDateTime                     types.String            `tfsdk:"quality_updates_rollback_start_date_time"`
	AllowWindows11Upgrade                                   types.Bool              `tfsdk:"allow_windows11_upgrade"`
	QualityUpdatesPaused                                    types.Bool              `tfsdk:"quality_updates_paused"`
	FeatureUpdatesPaused                                    types.Bool              `tfsdk:"feature_updates_paused"`
	SkipChecksBeforeRestart                                 types.Bool              `tfsdk:"skip_checks_before_restart"`
	BusinessReadyUpdatesOnly                                types.String            `tfsdk:"business_ready_updates_only"`
	AutomaticUpdateMode                                     types.String            `tfsdk:"automatic_update_mode"`
	DeliveryOptimizationMode                                types.String            `tfsdk:"delivery_optimization_mode"`
	PrereleaseFeatures                                      types.String            `tfsdk:"prerelease_features"`
	UpdateWeeks                                             types.String            `tfsdk:"update_weeks"`
	ActiveHoursStart                                        types.String            `tfsdk:"active_hours_start"`
	ActiveHoursEnd                                          types.String            `tfsdk:"active_hours_end"`
	UserPauseAccess                                         types.String            `tfsdk:"user_pause_access"`
	UserWindowsUpdateScanAccess                             types.String            `tfsdk:"user_windows_update_scan_access"`
	UpdateNotificationLevel                                 types.String            `tfsdk:"update_notification_level"`
	FeatureUpdatesRollbackWindowInDays                      types.Int32             `tfsdk:"feature_updates_rollback_window_in_days"`
	UninstallSettings                                       *UninstallSettingsModel `tfsdk:"uninstall"`
	UpdateActions                                           *UpdateActionsModel     `tfsdk:"update_actions"`
	DeadlineSettings                                        *DeadlineSettingsModel  `tfsdk:"deadline_settings"`
	EngagedRestartDeadlineInDays                            types.Int32             `tfsdk:"engaged_restart_deadline_in_days"`
	EngagedRestartSnoozeScheduleInDays                      types.Int32             `tfsdk:"engaged_restart_snooze_schedule_in_days"`
	EngagedRestartTransitionScheduleInDays                  types.Int32             `tfsdk:"engaged_restart_transition_schedule_in_days"`
	AutoRestartNotificationDismissal                        types.String            `tfsdk:"auto_restart_notification_dismissal"`
	ScheduleRestartWarningInHours                           types.Int32             `tfsdk:"schedule_restart_warning_in_hours"`
	ScheduleImminentRestartWarningInMinutes                 types.Int32             `tfsdk:"schedule_imminent_restart_warning_in_minutes"`
	EngagedRestartSnoozeScheduleForFeatureUpdatesInDays     types.Int32             `tfsdk:"engaged_restart_snooze_schedule_for_feature_updates_in_days"`
	EngagedRestartTransitionScheduleForFeatureUpdatesInDays types.Int32             `tfsdk:"engaged_restart_transition_schedule_for_feature_updates_in_days"`
	Assignments                                             types.Set               `tfsdk:"assignments"`
	Timeouts                                                timeouts.Value          `tfsdk:"timeouts"`
}

// UninstallSettingsModel defines the schema for update uninstall/rollback settings
type UninstallSettingsModel struct {
	FeatureUpdatesWillBeRolledBack types.Bool `tfsdk:"feature_updates_will_be_rolled_back"`
	QualityUpdatesWillBeRolledBack types.Bool `tfsdk:"quality_updates_will_be_rolled_back"`
}

// UpdateActionsModel defines the schema for update control actions
type UpdateActionsModel struct {
	FeatureUpdates *FeatureUpdateActionsModel `tfsdk:"feature_updates"`
	QualityUpdates *QualityUpdateActionsModel `tfsdk:"quality_updates"`
}

// FeatureUpdateActionsModel defines the schema for feature update actions
type FeatureUpdateActionsModel struct {
	Pause           types.Bool `tfsdk:"pause"`
	ExtendPause     types.Bool `tfsdk:"extend_pause"`
	TriggerUninstall types.Bool `tfsdk:"trigger_uninstall"`
}

// QualityUpdateActionsModel defines the schema for quality update actions
type QualityUpdateActionsModel struct {
	Pause           types.Bool `tfsdk:"pause"`
	TriggerUninstall types.Bool `tfsdk:"trigger_uninstall"`
}

// DeadlineSettingsModel defines the schema for deadline settings
type DeadlineSettingsModel struct {
	DeadlineForFeatureUpdatesInDays  types.Int32 `tfsdk:"deadline_for_feature_updates_in_days"`
	DeadlineForQualityUpdatesInDays  types.Int32 `tfsdk:"deadline_for_quality_updates_in_days"`
	DeadlineGracePeriodInDays        types.Int32 `tfsdk:"deadline_grace_period_in_days"`
	PostponeRebootUntilAfterDeadline types.Bool  `tfsdk:"postpone_reboot_until_after_deadline"`
}
